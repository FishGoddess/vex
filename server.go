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
	"sync"
	"syscall"
	"time"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")
	errCloseTimeout          = errors.New("vex: close server timeout")
)

// PacketHandler is a handler for handling packets.
// You will receive a byte slice of request and should return a byte slice or error if necessary.
type PacketHandler func(ctx context.Context, requestBody []byte) (responseBody []byte, err error)

// Server is the vex server.
type Server struct {
	config        config
	listener      net.Listener
	handlers      map[PacketType]PacketHandler
	eventListener EventListener
	lock          sync.RWMutex
}

// NewServer returns a new vex server.
func NewServer(network string, address string, opts ...Option) *Server {
	config := newDefaultConfig(network, address).ApplyOptions(opts)
	return &Server{
		config:        *config,
		handlers:      make(map[PacketType]PacketHandler, 16),
		eventListener: config.EventListener,
	}
}

// Name returns the name of the server.
func (s *Server) Name() string {
	return s.config.Name
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

// handleConn handles one conn.
func (s *Server) handleConn(conn net.Conn) {
	var err error

	reader := bufio.NewReaderSize(conn, int(s.config.ReadBufferSize))
	writer := bufio.NewWriterSize(conn, int(s.config.WriteBufferSize))
	defer func() {
		err := writer.Flush()
		if err != nil {
			log("vex: writer flushes failed [%+v]", err)
		}
	}()

	localAddr := conn.LocalAddr()
	remoteAddr := conn.RemoteAddr()
	s.eventListener.CallOnServerGotConnected(ServerGotConnectedEvent{Server: s, LocalAddr: localAddr, RemoteAddr: remoteAddr})
	defer s.eventListener.CallOnServerGotDisconnected(ServerGotDisconnectedEvent{Server: s, LocalAddr: localAddr, RemoteAddr: remoteAddr})

	ctx := wrapContext(context.Background(), localAddr, remoteAddr)
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
			log("vex: handler for %+v not found", packetType)
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
	s.eventListener.CallOnServerStart(ServerStartEvent{Server: s})
	defer s.eventListener.CallOnServerShutdown(ServerShutdownEvent{Server: s})

	var wg sync.WaitGroup
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// This error means listener has been closed.
			if errors.Is(err, net.ErrClosed) {
				break
			}

			log("vex: listener accepts failed %+v", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				err := conn.Close()
				if err != nil {
					log("vex: close connection failed [%+v]", err)
					return
				}
			}()

			err := conn.SetDeadline(time.Now().Add(s.config.ConnTimeout))
			if err != nil {
				log("vex: set deadline to connection failed [%+v]", err)
				return
			}

			s.handleConn(conn)
		}()
	}

	// Set a timer, so we won't wait too long.
	waitCh := make(chan struct{})
	go func() {
		wg.Wait() // Wait() won't stop if close is timeout and there are many connections waiting for handling.
		waitCh <- struct{}{}
	}()

	timer := time.NewTimer(s.config.CloseTimeout)
	defer timer.Stop()

	select {
	case <-waitCh:
		return nil
	case <-timer.C:
		return errCloseTimeout
	}
}

// ListenAndServe listens on address in network and begins serving.
func (s *Server) ListenAndServe() (err error) {
	s.listener, err = net.Listen(s.config.network, s.config.address)
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
