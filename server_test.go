// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
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

		sequence := uint64(12345)
		packet := packets.Packet{Magic: packets.Magic, Type: packets.PacketTypeRequest, Sequence: sequence}
		packet.With([]byte("test"))

		err = packets.Encode(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := strings.Repeat("test\n", i)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		decodePacket, err := packets.Decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if decodePacket.Magic != packets.Magic {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Magic, packets.Magic)
		}

		if decodePacket.Type != packets.PacketTypeResponse {
			t.Fatalf("%d: got %v != want %v", i, decodePacket.Type, packets.PacketTypeResponse)
		}

		if decodePacket.Length != uint32(len(decodePacket.Data)) {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Length, uint32(len(packet.Data)))
		}

		if decodePacket.Sequence != sequence {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Sequence, sequence)
		}

		got = string(decodePacket.Data)
		want = string(packet.Data)
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

		sequence := uint64(12345)
		packet := packets.Packet{Magic: packets.Magic, Type: 0, Sequence: sequence}
		packet.With([]byte("test"))

		err = packets.Encode(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := ""
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		decodePacket, err := packets.Decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if decodePacket.Magic != packets.Magic {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Magic, packets.Magic)
		}

		if decodePacket.Type != packets.PacketTypeError {
			t.Fatalf("%d: got %v != want %v", i, decodePacket.Type, packets.PacketTypeError)
		}

		if decodePacket.Length != uint32(len(decodePacket.Data)) {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Length, uint32(len(decodePacket.Data)))
		}

		if decodePacket.Sequence != sequence {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Sequence, sequence)
		}

		got = string(decodePacket.Data)
		want = "vex: packet type 0 is wrong"
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

		sequence := uint64(12345)
		packet := packets.Packet{Magic: packets.Magic, Type: packets.PacketTypeRequest, Sequence: sequence}
		packet.With([]byte("test"))

		err = packets.Encode(conn, packet)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)

		got := string(handler.data)
		want := strings.Repeat("test\n", i)
		if got != want {
			t.Fatalf("%d: got %s != want %s", i, got, want)
		}

		decodePacket, err := packets.Decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if decodePacket.Magic != packets.Magic {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Magic, packets.Magic)
		}

		if decodePacket.Type != packets.PacketTypeError {
			t.Fatalf("%d: got %v != want %v", i, decodePacket.Type, packets.PacketTypeError)
		}

		if decodePacket.Length != uint32(len(decodePacket.Data)) {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Length, uint32(len(decodePacket.Data)))
		}

		if decodePacket.Sequence != sequence {
			t.Fatalf("%d: got %d != want %d", i, decodePacket.Sequence, sequence)
		}

		got = string(decodePacket.Data)
		want = string(packet.Data)
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
