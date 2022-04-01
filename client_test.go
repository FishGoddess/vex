// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"net"
	"testing"
	"time"
)

func runTestServer(t *testing.T, address string, str string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	buffer := make([]byte, 64)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	magic := int32(endian.Uint32(buffer[:magicSize]))
	if magic != magicNumber {
		t.Errorf("magic %d != magicNumber %d", magic, magicNumber)
	}

	version := buffer[magicSize]
	if version != protocolVersion {
		t.Errorf("version %d != protocolVersion %d", version, protocolVersion)
	}

	packetType := buffer[magicSize+versionSize]
	if packetType != packetTypeTest {
		t.Errorf("packetType %d != packetTypeTest %d", packetType, packetTypeTest)
	}

	bodySize := endian.Uint32(buffer[magicSize+versionSize+typeSize : headerSize])
	requestBody := buffer[headerSize : headerSize+bodySize]
	if string(requestBody) != str {
		t.Errorf("requestBody %v is wrong!", string(requestBody))
	}

	responseBody := []byte(str)
	bodySize = uint32(len(responseBody))
	header := make([]byte, headerSize)
	endian.PutUint32(header[:magicSize], magicNumber)
	header[magicSize] = protocolVersion
	header[magicSize+versionSize] = packetTypeOK
	endian.PutUint32(header[magicSize+versionSize+typeSize:headerSize], bodySize)

	n, err = conn.Write(header)
	if err != nil {
		t.Error(err)
	}

	if n != headerSize {
		t.Errorf("n %d != headerSize %d", n, headerSize)
	}

	n, err = conn.Write(responseBody)
	if err != nil {
		t.Error(err)
	}

	if n != int(bodySize) {
		t.Errorf("n %d != bodySize %d", n, bodySize)
	}
}

// go test -v -cover -run=^TestClient$
func TestClient(t *testing.T) {
	address := "127.0.0.1:5837"
	str := "key value"
	go runTestServer(t, address, str)
	time.Sleep(time.Second)

	client, err := NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	rsp, err := client.Send(packetTypeTest, []byte(str))
	if err != nil {
		t.Error(err)
	}

	if string(rsp) != str {
		t.Errorf("body %s is wrong!", string(rsp))
	}
}
