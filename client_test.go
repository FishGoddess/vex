// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:00:26

package vex

import (
	"net"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewClient$
func TestNewClient(t *testing.T) {

	listener, err := net.Listen("tcp", ":5837")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}

		request := []byte{
			ProtocolVersion, 2, 0, 0, 0, 2, 0, 0, 0, 3, 'k', 'e', 'y', 0, 0, 0, 5, 'v', 'a', 'l', 'u', 'e',
		}
		buffer = buffer[:n]
		if string(buffer) != string(request) {
			t.Fatalf("Request %v is wrong!", string(buffer))
		}

		n, err = conn.Write([]byte{
			ProtocolVersion, SuccessReply, 0, 0, 0, 2, 'o', 'k',
		})
		if err != nil {
			t.Fatal(err)
		}

		if n != 8 {
			t.Fatalf("Written count %d is wrong!", n)
		}
	}()

	time.Sleep(time.Second)

	client, err := NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	body, err := client.Do(2, [][]byte{
		[]byte("key"), []byte("value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "ok" {
		t.Fatalf("Body %s is wrong!", string(body))
	}
}
