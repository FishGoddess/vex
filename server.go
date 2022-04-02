// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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

// handleConnOK handles ok happening on conn.
func (s *Server) handleConnOK(writer io.Writer, body []byte) {
	err := writePacket(writer, packetTypeOK, body)
	if err != nil {
		log("vex: write ok packet failed with err %+v", err)
	}
}

// handleConnErr handles errors happening on conn.
func (s *Server) handleConnErr(writer io.Writer, err error) {
	err = writePacket(writer, packetTypeErr, []byte(err.Error()))
	if err != nil {
		log("vex: write err packet failed with err %+v", err)
	}
}

// handleConn handles one conn.
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	notify(eventConnected)
	defer notify(eventDisconnected)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		if writer.Buffered() > 0 {
			err := writer.Flush()
			if err != nil {
				log("vex: writer flushes failed with err [%+v]", err)
			}
		}

		packetType, requestBody, err := readPacket(reader)
		if err != nil {
			if err != io.EOF {
				log("vex: read packet failed with err [%+v]", err)
				s.handleConnErr(writer, err)
			}
			return
		}

		s.lock.RLock()
		handle, ok := s.handlers[packetType]
		s.lock.RUnlock()

		if !ok {
			log("vex: handler of %+v not found", packetType)
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

// serve runs the accepting task.
func (s *Server) serve() error {
	notify(eventServing)
	defer notify(eventShutdown)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// This error means listener has been closed
			// See src/internal/poll/fd.go@ErrNetClosing
			// TODO So ugly...
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}

			log("vex: listener accepts failed with err %+v", err)
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

	go s.listenOnSignals()
	return s.serve()
}

// listenOnSignals listens on signal so server can respond to some signals.
func (s *Server) listenOnSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log("vex: received signal %+v...", sig)
	if err := s.Close(); err != nil {
		log("vex: server closes failed with err %+v", err)
	}
}

// Close closes current server.
func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
