// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 15:56:59

package vex

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// go test -v -cover -run=^TestReadResponseFrom$
func TestReadResponseFrom(t *testing.T) {

	response := bytes.NewBuffer([]byte{
		ProtocolVersion, SuccessReply, 0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o',
	})

	reply, body, err := readResponseFrom(response)
	if err != nil {
		t.Fatal(err)
	}

	if reply != SuccessReply {
		t.Fatalf("Reply %d is wrong!", reply)
	}

	if string(body) != "hello" {
		t.Fatalf("Body %s is wrong!", string(body))
	}
}

// go test -v -cover -run=^TestWriteResponseTo$
func TestWriteResponseTo(t *testing.T) {

	buffer := bytes.NewBuffer(make([]byte, 0, 64))

	body := []byte("hello")
	n, err := writeResponseTo(buffer, 1, body)
	if err != nil {
		t.Fatal(err)
	}

	if n != headerLengthInProtocol+len(body) {
		t.Fatalf("Written count %d is wrong!", n)
	}

	response := buffer.Bytes()
	if len(response) != headerLengthInProtocol+len(body) {
		t.Fatalf("Response length %d is wrong!", len(response))
	}

	if response[0] != ProtocolVersion {
		t.Fatalf("Version %d is wrong!", response[0])
	}

	if response[1] != 1 {
		t.Fatalf("Reply %d is wrong!", response[1])
	}

	bodyLengthBytes := response[2:6]
	bodyLength := binary.BigEndian.Uint32(bodyLengthBytes)
	if bodyLength != uint32(len(body)) {
		t.Fatalf("Body length %d is wrong!", bodyLength)
	}

	body = response[6:]
	if string(body) != "hello" {
		t.Fatalf("Body %s is wrong!", string(body))
	}
}
