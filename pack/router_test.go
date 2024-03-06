// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

var (
	errTestRequestFailed = errors.New("vex: test request failed")
)

func checkTestBytes(t *testing.T, buf []byte, expectedPacketType PacketType, expectedPacketData string) {
	//t.Log(buffer)
	if len(buf) < headerSize {
		t.Errorf("len(buf) %d < headerSize %d", len(buf), headerSize)
	}

	header := Endian.Uint64(buf[:headerSize])
	magic := (header >> (typeBits + dataSizeBits)) & maxMagic
	if magic != magicNumber {
		t.Errorf("magic %d != magicNumber %d", magic, magicNumber)
	}

	packetType := PacketType((header >> dataSizeBits) & maxType)
	if packetType != expectedPacketType {
		t.Errorf("packetType %d != expectedPacketType %d", packetType, expectedPacketType)
	}

	dataSize := int(header & maxDataSize)
	if dataSize != len([]byte(expectedPacketData)) {
		t.Errorf("dataSize %d != len([]byte(expectedPacketData)) %d ", dataSize, len([]byte(expectedPacketData)))
	}

	data := buf[headerSize : headerSize+dataSize]
	if string(data) != expectedPacketData {
		t.Errorf("data %s != expectedPacketData %s", data, expectedPacketData)
	}
}

func runTestClient(t *testing.T, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	// Err test
	_, err = conn.Write([]byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}

	msg := errTestRequestFailed.Error()
	buf := make([]byte, headerSize+len(msg))

	n, err := io.ReadFull(conn, buf)
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buf[:n], packetTypeError, msg)

	// OK test
	_, err = conn.Write([]byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 9, 'k', 'e', 'y', ' ', 'v', 'a', 'l', 'u', 'e'})
	if err != nil {
		t.Error(err)
	}

	msg = "key value"
	buf = make([]byte, headerSize+len(msg))

	n, err = conn.Read(buf[:])
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buf[:n], packetTypeStandard, msg)
}

// go test -v -cover -count=1 -test.cpu=1 -run=^TestRouterHandle$
func TestRouterHandle(t *testing.T) {
	address := "127.0.0.1:8899"

	router := NewRouter()
	router.Register(packetTypeTest, func(ctx context.Context, packetType PacketType, requestPacket []byte) (responsePacket []byte, err error) {
		if len(requestPacket) <= 0 {
			return nil, errTestRequestFailed
		}

		responsePacket = requestPacket
		return requestPacket, nil
	})

	server := vex.NewServer(address, router.Handle)
	defer server.Close()

	go func() {
		if err := server.Serve(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)
	runTestClient(t, address)
}
