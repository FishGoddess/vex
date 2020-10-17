// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 16:11:56

package vex

import (
	"errors"
	"net"
	"strings"
	"sync"
)

var (
	commandHandlerNotFoundErr = errors.New("failed to find a handler of command")
)

type Server struct {
	listener net.Listener
	handlers map[byte]func(args [][]byte) (body []byte, err error)
}

func NewServer() *Server {
	return &Server{
		handlers: map[byte]func(args [][]byte) (body []byte, err error){},
	}
}

func (s *Server) RegisterHandler(command byte, handler func(args [][]byte) (body []byte, err error)) {
	s.handlers[command] = handler
}

func (s *Server) ListenAndServe(network string, address string) (err error) {

	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
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

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.handleConn(conn)
		}()
	}

	wg.Wait()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {

	reader := conn
	defer conn.Close()

	for {
		command, args, err := readRequestFrom(reader)
		if err != nil {
			if err == ProtocolVersionMismatchErr {
				continue
			}
			return
		}

		reply, body, err := s.handleRequest(command, args)
		if err != nil {
			writeErrorResponseTo(conn, err.Error())
			continue
		}

		_, err = writeResponseTo(conn, reply, body)
		if err != nil {
			continue
		}
	}
}

func (s *Server) handleRequest(command byte, args [][]byte) (reply byte, body []byte, err error) {
	handle, ok := s.handlers[command]
	if !ok {
		return ErrorReply, nil, commandHandlerNotFoundErr
	}

	body, err = handle(args)
	if err != nil {
		return ErrorReply, body, err
	}
	return SuccessReply, body, err
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
