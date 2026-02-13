// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	packets "github.com/FishGoddess/vex/internal/packet"
)

type testHandler struct {
	data []byte
	lock sync.Mutex
}

func (h *testHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.data = append(h.data, data...)
	h.data = append(h.data, '\n')
	return data, nil
}

// go test -v -cover -run=^xxx$
func TestNewServer(t *testing.T) {
	handler := new(testHandler)

	t.Run("nil address", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("nil address returns a nil recover")
			}
		}()

		NewServer("", handler)
	})

	t.Run("nil handler", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("nil handler returns a nil recover")
			}
		}()

		NewServer("127.0.0.1:0", nil)
	})

	t.Run("wrong address", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("new server panics")
			}
		}()

		svr := NewServer("127.0.0.1", handler)

		if err := svr.Serve(); err == nil {
			t.Fatal(err)
		}

		if err := svr.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("new server", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("new server panics")
			}
		}()

		svr := NewServer("127.0.0.1:0", handler)

		go func() {
			if err := svr.Serve(); err != nil {
				t.Error(err)
			}
		}()

		time.Sleep(time.Second)

		err := syscall.Kill(os.Getpid(), syscall.SIGQUIT)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)

		select {
		case <-svr.(*server).ctx.Done():
			t.Log("server context is done")
		default:
			t.Fatal("server context not done")
		}

		if err = svr.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("double serve", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("new server panics")
			}
		}()

		svr := NewServer("127.0.0.1:0", handler)

		go func() {
			if err := svr.Serve(); err != nil {
				t.Error(err)
			}
		}()

		time.Sleep(time.Second)

		if err := svr.Serve(); err != errServerAlreadyServing {
			t.Fatal(err)
		}

		if err := svr.Close(); err != nil {
			t.Fatal(err)
		}
	})
}

// go test -v -cover -run=^TestServerHandler$
func TestServerHandler(t *testing.T) {
	handler := new(testHandler)
	svr := NewServer("127.0.0.1:0", handler)

	go func() {
		if err := svr.Serve(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	address := svr.(*server).listener.Addr().String()

	testCase := func(i int) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Fatal(err)
		}

		defer conn.Close()

		id := uint64(i)
		data := []byte("test")
		packet := packets.New(id)
		packet.SetData(data)

		err = packets.WritePacket(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := strings.Repeat("test\n", i)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		readPacket, err := packets.ReadPacket(conn)
		if err != nil {
			t.Fatal(err)
		}

		if readPacket.ID() != id {
			t.Fatalf("%d: got %d != want %d", i, readPacket.ID(), id)
		}

		readData, err := readPacket.Data()
		if err != nil {
			t.Fatal(err)
		}

		got = string(readData)
		want = string(data)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}
	}

	for i := range 10 {
		testCase(i + 1)
	}

	if err := svr.Close(); err != nil {
		t.Fatal(err)
	}
}

type testErrorHandler struct {
	data []byte
	lock sync.Mutex
}

func (h *testErrorHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.data = append(h.data, data...)
	h.data = append(h.data, '\n')
	return nil, errors.New(string(data))
}

// go test -v -cover -run=^TestServerError$
func TestServerError(t *testing.T) {
	handler := new(testErrorHandler)
	svr := NewServer("127.0.0.1:0", handler)

	go func() {
		if err := svr.Serve(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	address := svr.(*server).listener.Addr().String()

	testCase := func(i int) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Fatal(err)
		}

		defer conn.Close()

		id := uint64(i)
		data := []byte("test")
		packet := packets.New(id)
		packet.SetData(data)

		err = packets.WritePacket(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := strings.Repeat("test\n", i)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		readPacket, err := packets.ReadPacket(conn)
		if err != nil {
			t.Fatal(err)
		}

		if readPacket.ID() != id {
			t.Fatalf("%d: got %d != want %d", i, readPacket.ID(), id)
		}

		_, err = readPacket.Data()
		if err == nil {
			t.Fatalf("data returns nil error")
		}

		got = string(err.Error())
		want = string(data)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}
	}

	for i := range 10 {
		testCase(i + 1)
	}

	if err := svr.Close(); err != nil {
		t.Fatal(err)
	}
}

// go test -v -cover -run=^TestServerErrorHandler$
func TestServerErrorHandler(t *testing.T) {
	handler := new(testErrorHandler)
	svr := NewServer("127.0.0.1:0", handler)

	go func() {
		if err := svr.Serve(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	address := svr.(*server).listener.Addr().String()

	testCase := func(i int) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Fatal(err)
		}

		defer conn.Close()

		id := uint64(i)
		data := []byte("test")
		packet := packets.New(id)
		packet.SetData(data)

		err = packets.WritePacket(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := strings.Repeat("test\n", i)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		readPacket, err := packets.ReadPacket(conn)
		if err != nil {
			t.Fatal(err)
		}

		if readPacket.ID() != id {
			t.Fatalf("%d: got %d != want %d", i, readPacket.ID(), id)
		}

		_, err = readPacket.Data()
		if err == nil {
			t.Fatalf("data returns nil error")
		}

		got = string(err.Error())
		want = string(data)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}
	}

	for i := range 10 {
		testCase(i + 1)
	}

	if err := svr.Close(); err != nil {
		t.Fatal(err)
	}
}
