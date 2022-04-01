// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")

	Log = log.Printf // Logger outputs some messages.
)

type PacketHandler func(requestBody []byte) (responseBody []byte, err error)

type Server struct {
	listener net.Listener
	handlers map[PacketType]PacketHandler
	wg       sync.WaitGroup
	lock     sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[PacketType]PacketHandler, 16),
	}
}

func (s *Server) RegisterPacketHandler(packetType PacketType, handler PacketHandler) {
	s.lock.Lock()
	s.handlers[packetType] = handler
	s.lock.Unlock()
}

func (s *Server) handleConnError(writer io.Writer, err error) {
	Log("vex: read packet failed with err %+v", err)
	err = writePacket(writer, packetTypeErr, []byte(err.Error()))
	if err != nil {
		Log("vex: write packet failed with err %+v", err)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		if writer.Buffered() > 0 {
			writer.Flush()
		}

		packetType, requestBody, err := readPacket(reader)
		if err != nil {
			s.handleConnError(writer, err)
			return
		}

		s.lock.RLock()
		handle, ok := s.handlers[packetType]
		s.lock.RUnlock()

		if !ok {
			s.handleConnError(writer, errPacketHandlerNotFound)
			continue
		}

		responseBody, err := handle(requestBody)
		if err != nil {
			s.handleConnError(writer, err)
			continue
		}

		err = writePacket(writer, packetTypeOK, responseBody)
		if err != nil {
			s.handleConnError(writer, err)
		}
	}
}

func (s *Server) serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// This error means listener has been closed
			// See src/internal/poll/fd.go@ErrNetClosing
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConn(conn)
		}()
	}

	s.wg.Wait()
	return nil
}

func (s *Server) ListenAndServe(network string, address string) (err error) {
	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}
	return s.serve()
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
