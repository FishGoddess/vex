// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	packets "github.com/FishGoddess/vex/internal/packet"
)

var (
	errServerAlreadyServing = errors.New("vex: server is already serving")
)

// Handler is for handling the data from client and returns the new data or an error if failed.
type Handler interface {
	Handle(ctx *Context, data []byte) ([]byte, error)
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
	conns    map[uint64]net.Conn
	connID   uint64
	handler  Handler

	group sync.WaitGroup
	lock  sync.RWMutex
}

// NewServer creates a server with address and handler.
func NewServer(address string, handler Handler, opts ...Option) Server {
	conf := newConfig().apply(opts...)

	if address == "" {
		panic("vex: server address is nil")
	}

	if handler == nil {
		panic("vex: server handler is nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := new(server)
	server.conf = conf
	server.ctx = ctx
	server.cancel = cancel
	server.address = address
	server.conns = make(map[uint64]net.Conn, 64)
	server.connID = 0
	server.handler = handler

	go server.watchSignals()
	return server
}

func (s *server) watchSignals() {
	logger := s.conf.logger

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case sg := <-signalCh:
		logger.Info("received a signal", "signal", sg)

		if err := s.Close(); err != nil {
			logger.Error("close server failed", "err", err)
		}
	case <-s.ctx.Done():
		logger.Debug("server context is done")
	}
}

func (s *server) nextConnID() uint64 {
	s.connID++
	return s.connID
}

func (s *server) handlePacket(reader io.Reader, writer io.Writer) error {
	packet, err := packets.ReadPacket(reader)
	if err != nil {
		return err
	}

	data, err := packet.Data()
	if err != nil {
		return err
	}

	ctx := acquireContext(s.ctx)
	defer releaseContext(ctx)

	data, err = s.handler.Handle(ctx, data)
	if err != nil {
		packet.SetError(err)
	} else {
		packet.SetData(data)
	}

	err = packets.WritePacket(writer, packet)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) handleConn(conn net.Conn) {
	logger := s.conf.logger

	reader := bufio.NewReader(conn)
	for {
		err := s.handlePacket(reader, conn)
		if err == io.EOF {
			logger.Debug("handle packet eof", "err", err)
			return
		}

		if errors.Is(err, net.ErrClosed) {
			logger.Debug("handle packet closed", "err", err)
			return
		}

		if err != nil {
			logger.Error("handle packet failed", "err", err)
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

		s.lock.Lock()
		connID := s.nextConnID()
		s.conns[connID] = conn
		s.lock.Unlock()

		s.group.Go(func() {
			defer conn.Close()

			defer func() {
				s.lock.Lock()
				delete(s.conns, connID)
				s.lock.Unlock()
			}()

			logger.Info("handle conn start", "address", conn.RemoteAddr())
			s.handleConn(conn)
			logger.Info("handle conn end", "address", conn.RemoteAddr())
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

		logger.Error("server is already serving", "address", s.address)
		return errServerAlreadyServing
	}

	var lc net.ListenConfig

	listener, err := lc.Listen(s.ctx, "tcp", s.address)
	if err != nil {
		s.lock.Unlock()

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
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.lock.Unlock()

			return err
		}
	}

	for _, conn := range s.conns {
		if err := conn.Close(); err != nil {
			s.lock.Unlock()

			return err
		}
	}

	s.cancel()
	s.listener = nil
	s.conns = nil
	s.connID = 0
	s.lock.Unlock()
	s.group.Wait()
	return nil
}
