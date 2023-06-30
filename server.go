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

var (
	errCloseServerTimeout = errors.New("vex: close server timeout")
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
			log.Error(fmt.Errorf("%+v", r), "recover from handling")
		}
	}()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Error(err, "close connection failed")
		}
	}()

	if err := setupConn(&s.Config, conn); err != nil {
		log.Error(err, "setup connection failed")
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
				log.Info("server listener closed")
				break
			}

			log.Error(err, "listener accepts failed")
			continue
		}

		log.Debug("accepted from %s", conn.RemoteAddr())

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
		return errCloseServerTimeout
	}
}

func (s *server) monitorSignals() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log.Info("received signal %+v", sig)

	if err := s.Close(); err != nil {
		log.Error(err, "close server failed")
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

func (s *server) Close() error {
	return s.listener.Close()
}
