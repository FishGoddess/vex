// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	packets "github.com/FishGoddess/vex/internal/packet"
)

// Handler is for handling the data from client and returns the new data or an error if failed.
type Handler interface {
	Handle(ctx context.Context, data []byte) ([]byte, error)
}

// Server is the interface of vex server.
type Server interface {
	Serve() error
	Close() error
}

type server struct {
	conf *config

	ctx    context.Context
	cancel context.CancelFunc

	address  string
	listener net.Listener
	handler  Handler

	group sync.WaitGroup
	lock  sync.RWMutex
}

// NewServer creates a server with address and handler.
func NewServer(address string, handler Handler, opts ...Option) Server {
	conf := newConfig().apply(opts...)
	ctx, cancel := context.WithCancel(context.Background())

	if address == "" {
		panic("vex: server address is nil")
	}

	if handler == nil {
		panic("vex: server handler is nil")
	}

	server := &server{
		conf:    conf,
		ctx:     ctx,
		cancel:  cancel,
		address: address,
		handler: handler,
	}

	go server.watchSignals()
	return server
}

func (s *server) watchSignals() {
	logger := s.conf.logger

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sg := <-signalCh
	logger.Info("received a signal", "signal", sg)

	if err := s.Close(); err != nil {
		logger.Error("close server failed", "err", err)
	}
}

func (s *server) handlePacket(reader io.Reader, writer io.Writer) error {
	logger := s.conf.logger

	packet, err := packets.Decode(reader)
	if err != nil {
		logger.Error("decode packet failed", "err", err)
		return err
	}

	if packet.Type == packets.PacketTypeRequest {
		data, err := s.handler.Handle(s.ctx, packet.Data)
		if err == nil {
			packet.Type = packets.PacketTypeResponse
			packet.With(data)
		} else {
			packet.Type = packets.PacketTypeError
			packet.With(data)
		}
	} else {
		err = fmt.Errorf("vex: packet type %v is wrong", packet.Type)

		packet.Type = packets.PacketTypeError
		packet.With([]byte(err.Error()))
	}

	err = packets.Encode(writer, packet)
	if err != nil {
		logger.Error("encode packet failed", "err", err, "packet", packet)
		return err
	}

	return nil
}

func (s *server) handleConn(conn net.Conn) {
	logger := s.conf.logger
	logger.Debug("handle conn", "address", conn.RemoteAddr())
	defer logger.Debug("handle conn done", "address", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			writer.Flush()
		}

		if err := s.handlePacket(reader, writer); err != nil {
			return
		}
	}
}

func (s *server) serve() error {
	logger := s.conf.logger
	logger.Info("server is serving", "address", s.address)

	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			logger.Info("listener is closed", "address", s.address)
			break
		}

		if err != nil {
			logger.Error("accept conn failed", "err", err)
			continue
		}

		s.group.Go(func() {
			defer conn.Close()

			s.handleConn(conn)
		})
	}

	return nil
}

// Serve serves on address and returns an error if failed.
func (s *server) Serve() error {
	logger := s.conf.logger

	s.lock.Lock()
	if s.listener != nil {
		s.lock.Unlock()

		return errors.New("vex: server is already serving")
	}

	var lc net.ListenConfig
	listener, err := lc.Listen(s.ctx, "tcp", s.address)
	if err != nil {
		logger.Error("listen tcp failed", "err", err, "address", s.address)
		return err
	}

	s.listener = listener
	s.lock.Unlock()
	return s.serve()
}

// Close closes the server and returns an error if failed.
func (s *server) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
	}

	s.cancel()
	s.group.Wait()
	return nil
}
