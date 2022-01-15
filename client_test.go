// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2022/01/15 02:08:38

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

		bodyLength := binary.BigEndian.Uint32(buffer[versionLength+tagLength : headerLength])

		buffer = buffer[headerLength : headerLength+bodyLength]
		if string(buffer) != str {
			t.Errorf("request %v is wrong!", string(buffer))
		}

		body := []byte(str)
		bodyLength = uint32(len(body))
		header := make([]byte, headerLength)
		header[0] = ProtocolVersion
		header[1] = okTag
		binary.BigEndian.PutUint32(header[versionLength+tagLength:headerLength], bodyLength)
		n, err = conn.Write(header)
		if err != nil {
			t.Error(err)
		}

		if n != headerLength {
			t.Errorf("written count %d is wrong!", n)
		}

		n, err = conn.Write(body)
		if err != nil {
			t.Error(err)
		}

		if n != int(bodyLength) {
			t.Errorf("written count %d is wrong!", n)
		}
	}()

	time.Sleep(time.Second)

	client, err := NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	rsp, err := client.Do(testTag, []byte(str))
	if err != nil {
		t.Error(err)
	}

	if string(rsp) != str {
		t.Errorf("body %s is wrong!", string(rsp))
	}
}
