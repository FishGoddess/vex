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

var (
	errAcceptConnTimeout = errors.New("vex: accept conn timeout")
)

// HandleFunc is a function for handling connected context.
// You should design your own handler function for your server.
type HandleFunc func(ctx *Context)

type Status struct {
	// Connected is the count of connected connections.
	Connected uint64 `json:"connected"`
}

type Server interface {
	io.Closer

	Serve() error
	Status() Status
}

type server struct {
	Config

	listener net.Listener
	handle   HandleFunc

	conns    chan net.Conn
	status   Status
	contexts *sync.Pool

	closeCh chan struct{}
	wg      sync.WaitGroup
	lock    sync.RWMutex
}

// NewServer creates a new server serving on address.
// Handler is an interface of handling a connection.
func NewServer(address string, handle HandleFunc, opts ...Option) Server {
	conf := newServerConfig(address).ApplyOptions(opts)

	contexts := &sync.Pool{New: func() any {
		return new(Context)
	}}

	server := &server{
		Config:   *conf,
		handle:   handle,
		conns:    make(chan net.Conn, conf.maxConnections),
		contexts: contexts,
		closeCh:  make(chan struct{}),
	}

	go server.handleConns()
	log.Debug("server %s using config %+v", server.name, server.Config)

	return server
}

func (s *server) newContext(conn net.Conn) *Context {
	ctx := s.contexts.Get().(*Context)
	ctx.setup(conn)

	return ctx
}

func (s *server) freeContext(ctx *Context) {
	s.contexts.Put(ctx)
}

func (s *server) handleConn(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()

	ctx := s.newContext(conn)
	defer s.freeContext(ctx)

	defer func() {
		if err := ctx.finish(); err != nil {
			log.Error(err, "server %s finished %s failed", s.name, remoteAddr)
		}
	}()

	log.Debug("server %s handles %s begin", s.name, remoteAddr)
	defer log.Debug("server %s handles %s end", s.name, remoteAddr)

	s.beforeHandling(ctx)
	defer s.afterHandling(ctx)

	s.handle(ctx)
}

func (s *server) handleConns() {
	for conn := range s.conns {
		s.wg.Add(1)
		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					log.Error(fmt.Errorf("%+v", r), "server %s recovered from handling %s", s.name, conn.RemoteAddr())
				}
			}()

			defer func() {
				s.wg.Done()

				s.lock.Lock()
				s.status.Connected--
				s.lock.Unlock()
			}()

			s.handleConn(conn)
		}(conn)
	}
}

func (s *server) enqueueConn(conn net.Conn) (err error) {
	timer := time.NewTimer(s.connectTimeout)
	defer timer.Stop()

	select {
	case s.conns <- conn:
		break
	case <-s.closeCh:
		return net.ErrClosed
	case <-timer.C:
		return errAcceptConnTimeout
	}

	s.lock.Lock()
	s.status.Connected++
	s.lock.Unlock()

	return nil
}

func (s *server) acceptConn(conn net.Conn) {
	var err error

	defer func() {
		if err == nil {
			return
		}

		if err = conn.Close(); err != nil {
			log.Error(err, "server %s closes %s failed", s.name, conn.RemoteAddr())
		}
	}()

	if err = setupConn(&s.Config, conn); err != nil {
		log.Error(err, "server %s setups %s failed", s.name, conn.RemoteAddr())
		return
	}

	if err = s.enqueueConn(conn); err != nil {
		log.Error(err, "server %s enqueues %s failed", s.name, conn.RemoteAddr())
	}
}

func (s *server) wait() error {
	go func() {
		s.wg.Wait()

		close(s.closeCh)
	}()

	close(s.conns)
	s.afterServing(s.address)

	timer := time.NewTimer(s.closeTimeout)
	defer timer.Stop()

	select {
	case <-s.closeCh:
		s.afterClosing(s.address)

		log.Info("server %s closed", s.name)
		return nil
	case <-timer.C:
		return fmt.Errorf("vex: close server %s timeout", s.name)
	}
}

func (s *server) serve() error {
	for {
		conn, err := s.listener.Accept()
		if err == nil {
			s.acceptConn(conn)
			continue
		}

		// Listener has been closed.
		if errors.Is(err, net.ErrClosed) {
			break
		}

		log.Error(err, "server %s accepts failed", s.name)
	}

	return s.wait()
}

func (s *server) monitorSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log.Info("server %s received signal %+v", s.name, sig)

	if err := s.Close(); err != nil {
		log.Error(err, "close server %s failed", s.name)
	}
}

func (s *server) Serve() error {
	listener, err := net.Listen(network, s.address)
	if err != nil {
		return err
	}

	s.lock.Lock()
	if s.listener != nil {
		return fmt.Errorf("vex: server %s already serving", s.name)
	}

	s.listener = listener
	s.lock.Unlock()

	go s.monitorSignals()
	log.Info("server %s is serving on %s", s.name, s.address)

	s.beforeServing(s.address)
	return s.serve()
}

// Status returns the status of the server.
func (s *server) Status() Status {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.status
}

func (s *server) Close() error {
	s.beforeClosing(s.address)

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}
