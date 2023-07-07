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
	errWaitForContextTimeout = errors.New("vex: wait for context timeout")
)

// HandleFunc is a function for handling connected context.
// You should design your own handler function for your server.
type HandleFunc func(ctx *Context)

type Status struct {
	// Connected is the count of connected connections.
	Connected uint64 `json:"connected"`

	// Waiting is the count of requests waiting for connecting.
	Waiting uint64 `json:"waiting"`
}

type Server interface {
	io.Closer

	Serve() error
	Status() Status
}

type server struct {
	Config

	listener *net.TCPListener
	handle   HandleFunc
	contexts chan *Context
	status   Status

	lock sync.RWMutex
}

// NewServer creates a new server serving on address.
// Handler is an interface of handling a connection.
func NewServer(address string, handle HandleFunc, opts ...Option) Server {
	conf := newServerConfig(address).ApplyOptions(opts)

	server := &server{
		Config: *conf,
		handle: handle,
	}

	return server
}

func (s *server) useContext() (ctx *Context, err error) {
	s.lock.Lock()
	if s.status.Connected < uint64(s.maxConnections) {
		s.lock.Unlock()

		return new(Context), nil
	}

	s.status.Waiting++
	s.lock.Unlock()

	if s.connectTimeout > 0 {
		timer := time.NewTimer(s.connectTimeout)
		defer timer.Stop()

		select {
		case ctx = <-s.contexts:
			break
		case <-timer.C:
			err = errWaitForContextTimeout
		}
	} else {
		if ctx = <-s.contexts; ctx == nil {
			err = errWaitForContextTimeout
		}
	}

	// The waiting count should be decreased no matter we got a context or an error.
	s.lock.Lock()
	s.status.Waiting--
	s.lock.Unlock()

	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (s *server) freeContext(ctx *Context) {
	remoteAddr := ctx.RemoteAddr()

	s.lock.Lock()
	s.status.Connected--
	s.lock.Unlock()

	if err := ctx.finish(); err != nil {
		log.Error(err, "server %s finished %s failed", s.name, remoteAddr)
		return
	}

	select {
	case s.contexts <- ctx:
		log.Debug("server %s frees context %s", s.name, remoteAddr)
	default:
		log.Debug("server %s discards context %s", s.name, remoteAddr)
	}
}

func (s *server) setupContext(ctx *Context, conn *net.TCPConn) error {
	if err := setupConn(&s.Config, conn); err != nil {
		return err
	}

	ctx.setup(conn)
	return nil
}

func (s *server) accept() (*Context, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("%+v", r), "server %s recovered from accepting", s.name)
		}
	}()

	ctx, err := s.useContext()
	if err != nil {
		return nil, err
	}

	conn, err := s.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err = s.setupContext(ctx, conn); err != nil {
		return nil, err
	}

	s.lock.Lock()
	s.status.Connected++
	s.lock.Unlock()

	return ctx, nil
}

func (s *server) handleContext(ctx *Context) {
	remoteAddr := ctx.RemoteAddr()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("%+v", r), "server %s recovered from handling %s", s.name, remoteAddr)
		}
	}()

	log.Debug("server %s handles %s begin", s.name, remoteAddr)
	defer log.Debug("server %s handles %s end", s.name, remoteAddr)

	s.beforeHandling(ctx)
	defer s.afterHandling(ctx)

	s.handle(ctx)
}

func (s *server) serve() error {
	var wg sync.WaitGroup
	for {
		ctx, err := s.accept()
		if err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer s.freeContext(ctx)

				s.handleContext(ctx)
			}()

			continue
		}

		// Listener has been closed.
		if errors.Is(err, net.ErrClosed) {
			break
		}

		log.Error(err, "server %s accepts failed", s.name)
	}

	closeCh := make(chan struct{})
	go func() {
		wg.Wait()

		close(s.contexts)
		close(closeCh)
	}()

	timer := time.NewTimer(s.closeTimeout)
	defer timer.Stop()

	select {
	case <-closeCh:
		log.Info("server %s closed", s.name)
		return nil
	case <-timer.C:
		return fmt.Errorf("vex: close server %s timeout", s.name)
	}
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

func (s *server) Serve() (err error) {
	s.beforeServing(s.address)
	defer s.afterServing(s.address, err)

	address, err := net.ResolveTCPAddr(network, s.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(network, address)
	if err != nil {
		return err
	}

	s.lock.Lock()
	if s.listener != nil {
		return fmt.Errorf("vex: server %s already serving", s.name)
	}

	s.listener = listener
	s.contexts = make(chan *Context, s.maxConnections)
	s.lock.Unlock()

	go s.monitorSignals()
	log.Info("server %s is serving on %s", s.name, s.address)

	return s.serve()
}

// Status returns the status of the server.
func (s *server) Status() Status {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.status
}

func (s *server) Close() (err error) {
	s.beforeClosing(s.address)
	defer s.afterClosing(s.address, err)

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener == nil {
		return nil
	}

	defer func() {
		if err == nil {
			s.listener = nil
			s.status = Status{}
		}
	}()

	return s.listener.Close()
}
