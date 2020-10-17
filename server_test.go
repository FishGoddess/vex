// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:16:51

package vex

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"
)

const (
	testCommand = byte(1)
)

var (
	testArgumentErr = errors.New("test command needs more arguments")
)

// go test -v -cover -run=^TestNewServer$
func TestNewServer(t *testing.T) {

	server := NewServer()
	server.RegisterHandler(testCommand, func(args [][]byte) (reply byte, body []byte, err error) {
		if len(args) < 2 {
			return ErrorReply, nil, testArgumentErr
		}
		body = bytes.Join(args, []byte{' '})
		return SuccessReply, body, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// failed test
	_, err = conn.Write([]byte{
		ProtocolVersion, testCommand, 0, 0, 0, 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	if buffer[0] != ProtocolVersion {
		t.Fatalf("Protocol version %d is wrong!", buffer[0])
	}

	if buffer[1] != ErrorReply {
		t.Fatalf("Reply %d is wrong!", buffer[1])
	}

	bodyLength := binary.BigEndian.Uint32(buffer[2:6])
	if bodyLength != uint32(len([]byte(testArgumentErr.Error()))) {
		t.Fatalf("Body length %d is wrong!", bodyLength)
	}

	buffer = buffer[6:n]
	if string(buffer) != testArgumentErr.Error() {
		t.Fatalf("Body %s is wrong!", string(buffer))
	}

	// successful test
	_, err = conn.Write([]byte{
		ProtocolVersion, testCommand, 0, 0, 0, 2, 0, 0, 0, 3, 'k', 'e', 'y', 0, 0, 0, 5, 'v', 'a', 'l', 'u', 'e',
	})
	if err != nil {
		t.Fatal(err)
	}

	buffer = make([]byte, 1024)
	n, err = conn.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	if buffer[0] != ProtocolVersion {
		t.Fatalf("Protocol version %d is wrong!", buffer[0])
	}

	if buffer[1] != SuccessReply {
		t.Fatalf("Reply %d is wrong!", buffer[1])
	}

	bodyLength = binary.BigEndian.Uint32(buffer[2:6])
	if bodyLength != uint32(9) {
		t.Fatalf("Body length %d is wrong!", bodyLength)
	}

	buffer = buffer[6:n]
	if string(buffer) != "key value" {
		t.Fatalf("Body %s is wrong!", string(buffer))
	}
}
