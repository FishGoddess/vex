// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
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

// Handler is a handler for handling connection.
type Handler interface {
	// Handle handles a connection with reader and writer.
	// Some information can be fetched in context.
	Handle(ctx context.Context, reader io.Reader, writer io.Writer)
}

type Server interface {
	io.Closer

	Serve() error
}

type server struct {
	Config

	handler  Handler
	listener *net.TCPListener
}

// NewServer creates a new server serving on address.
// Handler is an interface of handling a connection.
func NewServer(address string, handler Handler, opts ...Option) Server {
	server := &server{
		Config:  *newServerConfig(address).ApplyOptions(opts),
		handler: handler,
	}

	return server
}

func (s *server) handleConn(conn *net.TCPConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("%+v", r), "server %s recovered from handling connection %s", s.Name, conn.RemoteAddr())
		}
	}()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Error(err, "server %s closes connection %s failed", s.Name, conn.RemoteAddr())
		}
	}()

	if err := setupConn(&s.Config, conn); err != nil {
		log.Error(err, "server %s setups connection %s failed", s.Name, conn.RemoteAddr())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.handler.Handle(ctx, conn, conn)
}

func (s *server) serve() error {
	var wg sync.WaitGroup
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			// Listener has been closed.
			if errors.Is(err, net.ErrClosed) {
				log.Debug("server %s listener closed", s.Name)
				break
			}

			log.Error(err, "server %s listener accepts failed", s.Name)
			continue
		}

		log.Debug("server %s accepted from %s", s.Name, conn.RemoteAddr())

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
	defer log.Info("server %s finished serving", s.Name)

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

func (s *server) Close() error {
	log.Debug("server %s is closing", s.Name)
	return s.listener.Close()
}
