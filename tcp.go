// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

const (
	networkTCP = "tcp"
)

var (
	errCloseTimeout = errors.New("vex: close server timeout")
)

type tcpServer struct {
	address  string
	listener *net.TCPListener
}

func NewTCP(address string) Server {
	return &tcpServer{
		address: address,
	}
}

func (ts *tcpServer) serve() error {
	var wg sync.WaitGroup
	for {
		conn, err := ts.listener.AcceptTCP()

		if err != nil {
			// This error means listener has been closed.
			if errors.Is(err, net.ErrClosed) {
				break
			}

			log.Printf("vex: listener accepts failed %+v", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Printf("vex: close connection failed [%+v]", err)
					return
				}
			}()

			err := conn.SetDeadline(time.Now().Add(time.Second))
			if err != nil {
				log.Printf("vex: set deadline to connection failed [%+v]", err)
				return
			}

			//s.handleConn(conn)
		}()
	}

	// Set a timer, so we won't wait too long.
	waitCh := make(chan struct{})
	go func() {
		wg.Wait() // Wait() won't stop if close is timeout and there are many connections waiting for handling.
		waitCh <- struct{}{}
	}()

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	select {
	case <-waitCh:
		return nil
	case <-timer.C:
		return errCloseTimeout
	}
}

func (ts *tcpServer) Serve() error {
	address, err := net.ResolveTCPAddr(networkTCP, ts.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(networkTCP, address)
	if err != nil {
		return err
	}

	ts.listener = listener
	go monitorSignals(ts)

	return ts.serve()
}

func (ts *tcpServer) Close() error {
	return ts.listener.Close()
}
