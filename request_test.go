// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 15:40:27

package vex

import (
	"bytes"
	"testing"
)

// go test -v -cover -run=^TestReadRequestFrom$
func TestReadRequestFrom(t *testing.T) {

	request := bytes.NewBuffer([]byte{
		ProtocolVersion, 1, 0, 0, 0, 2, 0, 0, 0, 1, 49, 0, 0, 0, 2, 50, 51,
	})

	command, args, err := readRequestFrom(request)
	if err != nil {
		t.Fatal(err)
	}

	if command != 1 {
		t.Fatalf("Command %d is wrong!", command)
	}

	if len(args) != 2 {
		t.Fatalf("Length of args %d is wrong!", len(args))
	}

	if args[0][0] != 49 || args[1][0] != 50 || args[1][1] != 51 {
		t.Fatalf("Args %v is wrong!", args)
	}
	t.Log(command, args)
}

// go test -v -cover -run=^TestWriteRequestTo$
func TestWriteRequestTo(t *testing.T) {

	buffer := bytes.NewBuffer(make([]byte, 0, 64))

	n, err := writeRequestTo(buffer, 1, [][]byte{
		[]byte("hello"), []byte("world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if n != 24 {
		t.Fatalf("Written count %d is wrong!", n)
	}

	request := []byte{
		ProtocolVersion, 1, 0, 0, 0, 2, 0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o', 0, 0, 0, 5, 'w', 'o', 'r', 'l', 'd',
	}
	if string(request) != string(buffer.Bytes()) {
		t.Fatalf("Request %s is wrong!", string(request))
	}
}
