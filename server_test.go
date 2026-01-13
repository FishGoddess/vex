// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	packets "github.com/FishGoddess/vex/internal/packet"
)

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %+v", msg, err)
	}
}

func assertEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v != want %v", msg, got, want)
	}
}

func assertContains(t *testing.T, got, substr, msg string) {
	t.Helper()
	if !strings.Contains(got, substr) {
		t.Fatalf("%s: got %q not contains %q", msg, got, substr)
	}
}

type testHandler struct {
	data []byte
}

func (h *testHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	h.data = append(h.data, data...)
	h.data = append(h.data, '\n')
	return h.data, nil
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
	want := "test\n"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	packet, err = packets.Decode(conn)
	if err != nil {
		t.Fatal(err)
	}

	if packet.Magic != packets.Magic {
		t.Fatalf("got %d != want %d", packet.Magic, packets.Magic)
	}

	if packet.Type != packets.PacketTypeResponse {
		t.Fatalf("got %v != want %v", packet.Type, packets.PacketTypeResponse)
	}

	if packet.Length != uint32(len(packet.Data)) {
		t.Fatalf("got %d != want %d", packet.Length, uint32(len(handler.data)))
	}

	if packet.Sequence != sequence {
		t.Fatalf("got %d != want %d", packet.Sequence, sequence)
	}

	got = string(packet.Data)
	want = "test\n"
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	// if err = svr.Close(); err != nil {
	// 	t.Fatal(err)
	// }
}
