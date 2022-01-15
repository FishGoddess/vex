// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2022/01/15 00:53:37

package vex

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"sync"
)

var (
	errHandlerNotFound = errors.New("vex: handler not found")
)

type Handler func(req []byte) (rsp []byte, err error)

type Server struct {
	listener net.Listener
	handlers map[Tag]Handler
	wg       sync.WaitGroup
	lock     sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[Tag]Handler, 16),
	}
}

func (s *Server) RegisterHandler(tag Tag, handler Handler) {
	s.lock.Lock()
	s.handlers[tag] = handler
	s.lock.Unlock()
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	for {
		tag, req, err := readFrom(reader)
		if err == errProtocolMismatch {
			continue
		}

		if err != nil {
			return
		}

		s.lock.RLock()
		handler, ok := s.handlers[tag]
		s.lock.RUnlock()

		if !ok {
			writeTo(writer, errTag, []byte(errHandlerNotFound.Error()))
			writer.Flush()
			continue
		}

		rsp, err := handler(req)
		if err != nil {
			writeTo(writer, errTag, []byte(err.Error()))
			writer.Flush()
			continue
		}

		writeTo(writer, okTag, rsp)
		writer.Flush()
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
