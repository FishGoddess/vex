// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/12 23:20:54

package vex

import (
	"io"
	"net"
	"strings"
	"sync"

	"github.com/FishGoddess/logit"
)

func init() {
	logit.Me().EnableFileInfo()
}

type Server struct {
	listener net.Listener
	handlers map[string]Handler
}

func NewServer() *Server {
	return &Server{
		handlers: map[string]Handler{},
	}
}

func (s *Server) RegisterHandler(command string, handler Handler) {
	s.handlers[command] = handler
}

func (s *Server) ListenAndServe(network string, address string) error {

	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.listener = listener

	connWg := &sync.WaitGroup{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			// The err will be "use of closed network connection" if listener has been closed.
			// Actually, this is a stupid way but em...
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}

			// Ignore other error
			continue
		}

		// Mark every connection.
		connWg.Add(1)
		go func(c net.Conn) {
			defer connWg.Done()
			defer conn.Close()
			s.handleConn(c)
		}(conn)
	}

	// Wait for all connections have done.
	connWg.Wait()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {

	for {
		req, err := readRequest(conn)
		if err != nil {
			if err != io.EOF {
				logit.Errorf("Failed to read request! Error is %s!", err.Error())
				readRequestErrorHandler(newContext(conn, req))
			}
			return
		}

		handler, ok := s.handlers[req.command]
		if !ok {
			logit.Errorf("Failed to find handler %s!", req.command)
			notFoundErrorHandler(newContext(conn, req))
			return
		}

		handler(newContext(conn, req))
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}
