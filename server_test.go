// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"errors"
	"testing"
	"time"
)

const (
	packetTypeTest PacketType = 1
)

var (
	errTestRequestFailed = errors.New("vex: test request failed")
)

func checkTestBytes(t *testing.T, buffer []byte, failedTest bool, expected string) {
	//t.Log(buffer)
	magic := int32(endian.Uint32(buffer[:magicSize]))
	if magic != magicNumber {
		t.Errorf("magic %d != magicNumber %d", magic, magicNumber)
	}

	version := buffer[magicSize]
	if version != protocolVersion {
		t.Errorf("version %d != protocolVersion %d", version, protocolVersion)
	}

	expectedPacketType := packetTypeOK
	if failedTest {
		expectedPacketType = packetTypeErr
	}

	packetType := buffer[magicSize+versionSize]
	if packetType != expectedPacketType {
		t.Errorf("packetType %d != expectedPacketType %d", packetType, expectedPacketType)
	}

	bodySize := endian.Uint32(buffer[magicSize+versionSize+typeSize : headerSize])
	if bodySize != uint32(len([]byte(expected))) {
		t.Errorf("bodySize %d != len([]byte(expected)) %d ", bodySize, len([]byte(expected)))
	}

	requestBody := buffer[headerSize : headerSize+bodySize]
	if string(requestBody) != expected {
		t.Errorf("requestBody %s != expected %s", string(requestBody), expected)
	}
}

func runTestClient(t *testing.T, address string) {
	conn, err := dial("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	// Err test
	_, err = conn.Write([]byte{
		0x7, 0x55, 0xDD, 0x8C, protocolVersion, packetTypeTest, 0, 0, 0, 0,
	})
	if err != nil {
		t.Error(err)
	}

	buffer := make([]byte, 64)
	_, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buffer, true, errTestRequestFailed.Error())

	// OK test
	_, err = conn.Write([]byte{
		0x7, 0x55, 0xDD, 0x8C, protocolVersion, packetTypeTest, 0, 0, 0, 9, 'k', 'e', 'y', ' ', 'v', 'a', 'l', 'u', 'e',
	})
	if err != nil {
		t.Error(err)
	}

	buffer = make([]byte, 64)
	_, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	checkTestBytes(t, buffer, false, "key value")
}

// go test -v -cover -run=^TestNewServer$
func TestNewServer(t *testing.T) {
	address := "127.0.0.1:5837"

	server := NewServer()
	server.RegisterPacketHandler(packetTypeTest, func(requestBody []byte) (responseBody []byte, err error) {
		if len(requestBody) <= 0 {
			return nil, errTestRequestFailed
		}

		responseBody = requestBody
		return responseBody, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", address)
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)
	runTestClient(t, address)
}
