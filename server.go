// Copyright 2022 FishGoddess.  All rights reserved.
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
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")
)

// PacketHandler is a handler for handling packets.
// You will receive a byte slice of request and should return a byte slice or error if necessary.
type PacketHandler func(ctx context.Context, requestBody []byte) (responseBody []byte, err error)

// Server is the vex server.
type Server struct {
	config       config
	listener     net.Listener
	handlers     map[PacketType]PacketHandler
	eventHandler EventHandler
	wg           sync.WaitGroup
	lock         sync.RWMutex
}

// NewServer returns a new vex server.
func NewServer(opts ...Option) *Server {
	config := newDefaultConfig().ApplyOptions(opts)
	return &Server{
		config:       *config,
		handlers:     make(map[PacketType]PacketHandler, 16),
		eventHandler: config.EventHandler,
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
		log("vex: write ok packet failed %+v", err)
	}
}

// handleConnErr handles errors happening on conn.
func (s *Server) handleConnErr(writer io.Writer, err error) {
	err = writePacket(writer, packetTypeErr, []byte(err.Error()))
	if err != nil {
		log("vex: write err packet failed %+v", err)
	}
}

// publishEvent publishes an event and gives it to event handler.
func (s *Server) publishEvent(ctx context.Context, e Event) {
	if s.eventHandler != nil {
		s.eventHandler.HandleEvent(ctx, e)
	}
}

// setupConn setups conn and returns an error if failed.
func (s *Server) setupConn(conn net.Conn) error {
	return conn.SetDeadline(time.Now().Add(s.config.ConnTimeout))
}

// newContext returns a context with conn.
func (s *Server) newContext(conn net.Conn) context.Context {
	return wrapContext(context.Background(), conn)
}

// handleConn handles one conn.
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	err := s.setupConn(conn)
	if err != nil {
		log("vex: set up connection failed [%+v]", err)
		return
	}

	ctx := s.newContext(conn)
	s.publishEvent(ctx, eventConnected)
	defer s.publishEvent(ctx, eventDisconnected)

	reader := bufio.NewReaderSize(conn, int(s.config.ReadBufferSize))
	writer := bufio.NewWriterSize(conn, int(s.config.WriteBufferSize))
	defer func() {
		err = writer.Flush()
		if err != nil {
			log("vex: writer flushes failed [%+v]", err)
		}
	}()

	for {
		if writer.Buffered() > 0 {
			err = writer.Flush()
			if err != nil {
				log("vex: writer flushes failed [%+v]", err)
			}
		}

		packetType, requestBody, err := readPacket(reader)
		if err != nil {
			if err != io.EOF {
				log("vex: read packet failed [%+v]", err)
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

		responseBody, err := handle(ctx, requestBody)
		if err != nil {
			s.handleConnErr(writer, err)
			continue
		}

		s.handleConnOK(writer, responseBody)
	}
}

// serve runs the accepting task.
func (s *Server) serve() error {
	ctx := context.Background()
	s.publishEvent(ctx, eventServing)
	defer s.publishEvent(ctx, eventShutdown)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// This error means listener has been closed.
			// See src/internal/poll/fd.go@ErrNetClosing.
			// So ugly...
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}

			log("vex: listener accepts failed %+v", err)
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

	go s.listenToSignals()
	return s.serve()
}

// listenToSignals listens to signals so server can respond to these signals.
func (s *Server) listenToSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log("vex: received signal %+v...", sig)
	if err := s.Close(); err != nil {
		log("vex: server closes failed %+v", err)
	}
}

// Close closes current server.
func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
