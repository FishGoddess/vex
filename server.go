// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	network = "tcp"
)

var (
	errCloseTimeout = errors.New("vex: close server timeout")
)

type Server struct {
	ServerConfig

	listener *net.TCPListener
	handler  Handler
}

func NewServer(address string, opts ...ServerOption) *Server {
	conf := newServerConfig(network, address).ApplyOptions(opts)

	return &Server{
		ServerConfig: *conf,
	}
}

func (s *Server) Handle(handler Handler) {
	s.handler = handler
}

func (s *Server) handleConn(conn *net.TCPConn) {
	now := time.Now()

	readDeadline := now.Add(s.ReadTimeout)
	writeDeadline := now.Add(s.WriteTimeout)

	if err := conn.SetReadDeadline(readDeadline); err != nil {
		logError(err, "set read deadline to connection failed")
		return
	}

	if err := conn.SetWriteDeadline(writeDeadline); err != nil {
		logError(err, "set write deadline to connection failed")
		return
	}

	if err := conn.SetReadBuffer(s.ReadBufferSize); err != nil {
		logError(err, "set read buffer size of connection failed")
		return
	}

	if err := conn.SetWriteBuffer(s.WriteBufferSize); err != nil {
		logError(err, "set write buffer size of connection failed")
		return
	}

	s.handler.Handle(newConnection(conn, s.ReadBufferSize, s.WriteBufferSize))
}

func (s *Server) serve() error {
	var wg sync.WaitGroup
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			// Listener has been closed.
			if errors.Is(err, net.ErrClosed) {
				logDebug("server listener closed")
				break
			}

			logError(err, "listener accepts failed")
			continue
		}

		wg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logError(fmt.Errorf("%+v", r), "recover from connection")
				}
			}()

			defer wg.Done()
			defer func() {
				if err := conn.Close(); err != nil {
					logError(err, "close connection failed")
				}
			}()

			s.handleConn(conn)
		}()
	}

	// Set a timer, so we won't wait too long.
	closeCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(closeCh)
	}()

	timer := time.NewTimer(s.CloseTimeout)
	defer timer.Stop()

	select {
	case <-closeCh:
		return nil
	case <-timer.C:
		return errCloseTimeout
	}
}

func (s *Server) monitorSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	logInfo("received signal %+v", sig)

	if err := s.Close(); err != nil {
		logError(err, "close server failed")
	}
}

func (s *Server) Serve() error {
	address, err := net.ResolveTCPAddr(network, s.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(network, address)
	if err != nil {
		return err
	}

	s.listener = listener
	go s.monitorSignals()

	logInfo("server %s is serving on %s", s.Name, s.address)
	return s.serve()
}

func (s *Server) Close() error {
	return s.listener.Close()
}
