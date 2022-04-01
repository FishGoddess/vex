// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"runtime"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewClientPool$
func TestNewClientPool(t *testing.T) {
	server := NewServer()
	server.RegisterPacketHandler(1, func(req []byte) (rsp []byte, err error) {
		return []byte("test"), nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", "127.0.0.1:5837")
		if err != nil {
			panic(err)
		}
	}()

	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	pool, err := NewClientPool(64, func() (Client, error) {
		return NewClient("tcp", "127.0.0.1:5837")
	})
	if err != nil {
		t.Fatal("new client pool failed", err)
	}
	defer pool.Close()

	for i := 0; i < 512; i++ {
		client := pool.Get()

		response, err := client.Send(1, []byte("test"))
		if err != nil {
			client.Close()
			t.Fatalf("do command %d failed with %+v", i, err)
		}

		if string(response) != "test" {
			client.Close()
			t.Fatalf("response %s is wrong", string(response))
		}

		client.Close()
	}
}
