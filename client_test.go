// Copyright 2025 FishGoddess. All rights reserved.
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

func runTestServer(t *testing.T, address string, ch chan struct{}) {
	listener, err := net.Listen(network, address)
	if err != nil {
		t.Error(err)
	}

	defer listener.Close()
	close(ch)

	conn, err := listener.Accept()
	if err != nil {
		t.Error(err)
	}

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[:])
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Error(err)
		}

		_, err = conn.Write(buf[:n])
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Error(err)
		}
	}
}

// go test -v -cover -run=^TestClient$
func TestClient(t *testing.T) {
	address := "127.0.0.1:12345"

	ch := make(chan struct{})
	go runTestServer(t, address, ch)

	<-ch
	time.Sleep(time.Second)

	client, err := NewClient(address)
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	var buf [1024]byte
	for i := 0; i < 100; i++ {
		msg := strconv.FormatUint(rand.Uint64(), 10)
		t.Log(msg)

		_, err = client.Write([]byte(msg))
		if err != nil {
			t.Error(err)
		}

		n, err := client.Read(buf[:])
		if err != nil {
			t.Error(err)
		}

		received := string(buf[:n])
		if received != msg {
			t.Errorf("received %s != msg %s", received, msg)
		}
	}
}
