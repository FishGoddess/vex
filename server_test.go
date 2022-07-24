// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"errors"
	"testing"
	"time"
)

var (
	errTestRequestFailed = errors.New("vex: test request failed")
)

func checkTestBytes(t *testing.T, buffer []byte, expectedPacketType PacketType, expectedBody string) {
	//t.Log(buffer)
	if len(buffer) < headerSize {
		t.Errorf("len(buffer) %d < headerSize %d", len(buffer), headerSize)
	}

	header := endian.Uint64(buffer[:headerSize])
	magic := (header >> (typeBits + bodySizeBits)) & maxMagic
	if magic != magicNumber {
		t.Errorf("magic %d != magicNumber %d", magic, magicNumber)
	}

	packetType := PacketType((header >> bodySizeBits) & maxType)
	if packetType != expectedPacketType {
		t.Errorf("packetType %d != expectedPacketType %d", packetType, expectedPacketType)
	}

	bodySize := int(header & maxBodySize)
	if bodySize != len([]byte(expectedBody)) {
		t.Errorf("bodySize %d != len([]byte(expectedBody)) %d ", bodySize, len([]byte(expectedBody)))
	}

	requestBody := buffer[headerSize : headerSize+bodySize]
	if string(requestBody) != expectedBody {
		t.Errorf("requestBody %s != expectedBody %s", string(requestBody), expectedBody)
	}
}

func runTestClient(t *testing.T, address string) {
	conn, err := dial("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	// Err test
	_, err = conn.Write([]byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 0})
	if err != nil {
		t.Error(err)
	}

	buffer := make([]byte, 64)
	_, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buffer, packetTypeErr, errTestRequestFailed.Error())

	// OK test
	_, err = conn.Write([]byte{0xC, 0x63, 0x8B, packetTypeTest, 0, 0, 0, 9, 'k', 'e', 'y', ' ', 'v', 'a', 'l', 'u', 'e'})
	if err != nil {
		t.Error(err)
	}

	buffer = make([]byte, 64)
	_, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buffer, packetTypeOK, "key value")
}

// go test -v -cover -run=^TestNewServer$
func TestNewServer(t *testing.T) {
	address := "127.0.0.1:5837"

	server := NewServer("tcp", address)
	server.RegisterPacketHandler(packetTypeTest, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		if len(requestBody) <= 0 {
			return nil, errTestRequestFailed
		}

		responseBody = requestBody
		return responseBody, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)
	runTestClient(t, address)
}
