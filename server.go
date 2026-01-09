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

type Server interface {
	Serve() error
	Close() error
}

type server struct {
	ctx    context.Context
	cancel context.CancelFunc

	address  string
	listener net.Listener

	group sync.WaitGroup
	lock  sync.RWMutex
}

func NewServer(address string) Server {
	ctx, cancel := context.WithCancel(context.Background())

	server := &server{
		ctx:     ctx,
		cancel:  cancel,
		address: address,
	}

	go server.watchSignals()
	return server
}

func (s *server) watchSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	_ = <-signalCh

	if err := s.Close(); err != nil {

	}
}

func (s *server) handlePacket(reader io.Reader, writer io.Writer) {
	packet, err := packets.Decode(reader)
	if err != nil {
		return
	}

	if packet.Type() != packets.PacketTypeRequest {
		return
	}

	s.group.Go(func() {
		// TODO 处理包
		packet.SetType(packets.PacketTypeResponse)

		err = packets.Encode(writer, packet)
		if err != nil {
			return
		}
	})
}

func (s *server) handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.handlePacket(reader, writer)
		}
	}
}

func (s *server) serve() error {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		}

		if err != nil {
			continue
		}

		s.group.Go(func() {
			defer conn.Close()

			s.handleConn(conn)
		})
	}

	return nil
}

func (s *server) Serve() error {
	var lc net.ListenConfig

	listener, err := lc.Listen(s.ctx, "tcp", s.address)
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.listener = listener
	s.lock.Unlock()

	return s.serve()
}

func (s *server) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener != nil {
		return s.listener.Close()
	}

	s.cancel()
	s.group.Wait()
	return nil
}
