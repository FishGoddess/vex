// Copyright 2025 FishGoddess. All rights reserved.
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
	errAcquireTokenTimeout = errors.New("vex: acquire token timeout")
)

// HandleFunc is a function for handling connected context.
// You should design your own handler function for your server.
type HandleFunc func(ctx *Context)

type Status struct {
	// Connected is the quantity of connected connections.
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

	status   Status
	contexts *sync.Pool

	tokenCh chan struct{}
	closeCh chan struct{}

	wg   sync.WaitGroup
	lock sync.RWMutex
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
		contexts: contexts,
		tokenCh:  make(chan struct{}, conf.maxConnections),
		closeCh:  make(chan struct{}),
	}

	return server
}

func (s *server) acquireToken() (shouldRetry bool, err error) {
	timer := time.NewTimer(s.connectTimeout)
	defer timer.Stop()

	select {
	case s.tokenCh <- struct{}{}:
		return true, nil
	case <-s.closeCh:
		return false, nil
	case <-timer.C:
		return false, errAcquireTokenTimeout
	}
}

func (s *server) releaseToken() {
	select {
	case <-s.tokenCh:
		return
	case <-s.closeCh:
		return
	}
}

func (s *server) newContext(conn net.Conn) *Context {
	ctx := s.contexts.Get().(*Context)
	ctx.setup(conn)

	return ctx
}

func (s *server) freeContext(ctx *Context) {
	if err := ctx.finish(); err != nil {
		log.Error(err, "server %s finished %s failed", s.name, ctx.RemoteAddr())
	}

	s.contexts.Put(ctx)
}

func (s *server) handleConn(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("%+v", r), "server %s recovered from handling %s", s.name, remoteAddr)
		}
	}()

	defer func() {
		s.lock.Lock()
		s.status.Connected--
		s.lock.Unlock()
	}()

	ctx := s.newContext(conn)
	defer s.freeContext(ctx)

	log.Debug("server %s handles %s begin", s.name, remoteAddr)
	defer log.Debug("server %s handles %s end", s.name, remoteAddr)

	s.beforeHandling(ctx)
	defer s.afterHandling(ctx)

	select {
	case <-ctx.Done():
		return
	default:
		s.handle(ctx)
	}
}

func (s *server) acceptConn(conn net.Conn) {
	var err error

	defer func() {
		if err == nil {
			return
		}

		s.releaseToken()
		log.Error(err, "server %s accepts %s failed", s.name, conn.RemoteAddr())

		if err = conn.Close(); err != nil {
			log.Error(err, "server %s closes %s failed", s.name, conn.RemoteAddr())
		}
	}()

	if err = setupConn(&s.Config, conn); err != nil {
		log.Error(err, "server %s setups %s failed", s.name, conn.RemoteAddr())
		return
	}

	s.lock.Lock()
	s.status.Connected++
	s.lock.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.releaseToken()

		s.handleConn(conn)
	}()
}

func (s *server) wait() error {
	go func() {
		s.wg.Wait()

		close(s.closeCh)
	}()

	close(s.tokenCh)
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
		shouldRetry, err := s.acquireToken()
		if err != nil {
			log.Error(err, "server %s acquires token failed", s.name)
			continue
		}

		if !shouldRetry {
			break
		}

		conn, err := s.listener.Accept()
		if err == nil {
			s.acceptConn(conn)
			continue
		}

		// Listener has been closed.
		if errors.Is(err, net.ErrClosed) {
			break
		}

		s.releaseToken()
		log.Error(err, "server %s accepts failed", s.name)
	}

	return s.wait()
}

func (s *server) monitorSignals() {
	signalCh := make(chan os.Signal, 1)
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
