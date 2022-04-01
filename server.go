// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")
)

// PacketHandler is a handler for handling packets.
// You will receive a byte slice of request and should return a byte slice or error if necessary.
type PacketHandler func(requestBody []byte) (responseBody []byte, err error)

// Server is the vex server.
type Server struct {
	listener net.Listener
	handlers map[PacketType]PacketHandler
	wg       sync.WaitGroup
	lock     sync.RWMutex
}

// NewServer returns a new vex server.
func NewServer() *Server {
	return &Server{
		handlers: make(map[PacketType]PacketHandler, 16),
	}
}

// RegisterPacketHandler registers handler of packetType.
func (s *Server) RegisterPacketHandler(packetType PacketType, handler PacketHandler) {
	s.lock.Lock()
	s.handlers[packetType] = handler
	s.lock.Unlock()
}

func (s *Server) handleConnOK(writer io.Writer, body []byte) {
	err := writePacket(writer, packetTypeOK, body)
	if err != nil {
		Log("vex: write ok packet failed with err %+v", err)
	}
}

func (s *Server) handleConnErr(writer io.Writer, err error) {
	err = writePacket(writer, packetTypeErr, []byte(err.Error()))
	if err != nil {
		Log("vex: write err packet failed with err %+v", err)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		if writer.Buffered() > 0 {
			err := writer.Flush()
			if err != nil {
				Log("vex: writer flushes failed with err [%+v]", err)
			}
		}

		packetType, requestBody, err := readPacket(reader)
		if err != nil {
			if err != io.EOF {
				Log("vex: read packet failed with err [%+v]", err)
				s.handleConnErr(writer, err)
			}
			return
		}

		s.lock.RLock()
		handle, ok := s.handlers[packetType]
		s.lock.RUnlock()

		if !ok {
			Log("vex: handler of %+v not found", packetType)
			s.handleConnErr(writer, errPacketHandlerNotFound)
			continue
		}

		responseBody, err := handle(requestBody)
		if err != nil {
			s.handleConnErr(writer, err)
			continue
		}

		s.handleConnOK(writer, responseBody)
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

// ListenAndServe listens on address in network and begins serving.
func (s *Server) ListenAndServe(network string, address string) (err error) {
	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}
	return s.serve()
}

// Close closes current server.
func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
