// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"io"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func testHandle(ctx *Context) {
	var buf [1024]byte

	for {
		n, err := ctx.Read(buf[:])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		_, err = ctx.Write(buf[:n])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}
}

func runTestClient(t *testing.T, address string, ch chan struct{}, closeCh chan struct{}) {
	<-ch
	time.Sleep(time.Second)

	conn, err := net.Dial(network, address)
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	var buf [1024]byte
	for i := 0; i < 100; i++ {
		msg := strconv.FormatUint(rand.Uint64(), 10)
		t.Log(msg)

		_, err := conn.Write([]byte(msg))
		if err != nil {
			t.Error(err)
		}

		n, err := conn.Read(buf[:])
		if err != nil {
			t.Error(err)
		}

		received := string(buf[:n])
		if received != msg {
			t.Errorf("received %s != msg %s", received, msg)
		}
	}

	close(closeCh)
}

// go test -v -cover -run=^TestServer$
func TestServer(t *testing.T) {
	address := "127.0.0.1:54321"

	ch := make(chan struct{})
	closeCh := make(chan struct{})
	go runTestClient(t, address, ch, closeCh)

	server := NewServer(address, testHandle)
	close(ch)

	go func() {
		<-closeCh
		server.Close()
	}()

	if err := server.Serve(); err != nil {
		t.Error(err)
	}
}
