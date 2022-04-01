// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewClient$
func TestNewClient(t *testing.T) {
	address := "127.0.0.1:5837"
	str := "key value"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		buffer := make([]byte, 128)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Error(err)
		}

		bodySize := binary.BigEndian.Uint32(buffer[versionSize+typeSize : headerSize])

		buffer = buffer[headerSize : headerSize+bodySize]
		if string(buffer) != str {
			t.Errorf("request %v is wrong!", string(buffer))
		}

		body := []byte(str)
		bodySize = uint32(len(body))
		header := make([]byte, headerSize)
		header[0] = protocolVersion
		header[1] = packetTypeTest
		binary.BigEndian.PutUint32(header[versionSize+typeSize:headerSize], bodySize)
		n, err = conn.Write(header)
		if err != nil {
			t.Error(err)
		}

		if n != headerSize {
			t.Errorf("written count %d is wrong!", n)
		}

		n, err = conn.Write(body)
		if err != nil {
			t.Error(err)
		}

		if n != int(bodySize) {
			t.Errorf("written count %d is wrong!", n)
		}
	}()

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
