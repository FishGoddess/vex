// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/FishGoddess/vex/log"
)

const (
	network = "tcp"
)

type HandleFunc func(ctx *Context)

type Server interface {
	io.Closer

	Serve() error
}

type server struct {
	Config

	handle   HandleFunc
	listener *net.TCPListener

	contextPool *sync.Pool
}

// NewServer creates a new server serving on address.
// Handler is an interface of handling a connection.
func NewServer(address string, handle HandleFunc, opts ...Option) Server {
	conf := newServerConfig(address).ApplyOptions(opts)

	contextPool := &sync.Pool{New: func() any {
		return new(Context)
	}}

	server := &server{
		Config:      *conf,
		handle:      handle,
		contextPool: contextPool,
	}

	return server
}

func (s *server) newContext(conn *net.TCPConn) *Context {
	ctx := s.contextPool.Get().(*Context)
	ctx.setup(conn)
	return ctx
}

func (s *server) releaseContext(ctx *Context) {
	s.contextPool.Put(ctx)
}

func (s *server) handleConn(conn *net.TCPConn) {
	remoteAddr := conn.RemoteAddr()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("%+v", r), "server %s recovered from handling connection %s", s.Name, remoteAddr)
		}
	}()

	if err := setupConn(&s.Config, conn); err != nil {
		log.Error(err, "server %s setups connection %s failed", s.Name, remoteAddr)
		return
	}

	ctx := s.newContext(conn)

	defer func() {
		if err := ctx.finish(); err != nil {
			log.Error(err, "server %s finished connection %s failed", s.Name, remoteAddr)
		}

		s.releaseContext(ctx)
	}()

	log.Debug("server %s handles connection %s begin", s.Name, remoteAddr)
	s.handle(ctx)
	log.Debug("server %s handles connection %s end", s.Name, remoteAddr)
}

func (s *server) serve() error {
	var wg sync.WaitGroup
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			// Listener has been closed.
			if errors.Is(err, net.ErrClosed) {
				log.Debug("server %s stopped listening", s.Name)
				break
			}

			log.Error(err, "server %s accepted failed", s.Name)
			continue
		}

		log.Debug("server %s accepts from %s", s.Name, conn.RemoteAddr())

		wg.Add(1)
		go func() {
			defer wg.Done()

			s.handleConn(conn)
		}()
	}

	closeCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(closeCh)
	}()

	timer := time.NewTimer(s.CloseTimeout)
	defer timer.Stop()

	select {
	case <-closeCh:
		log.Info("server %s closed", s.Name)
		return nil
	case <-timer.C:
		return fmt.Errorf("vex: close server %s timeout", s.Name)
	}
}

func (s *server) monitorSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log.Info("server %s received signal %+v", s.Name, sig)

	if err := s.Close(); err != nil {
		log.Error(err, "close server %s failed", s.Name)
	}
}

func (s *server) Serve() error {
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

	log.Info("server %s is serving on %s", s.Name, s.address)
	return s.serve()
}

func (s *server) Close() (err error) {
	return s.listener.Close()
}
