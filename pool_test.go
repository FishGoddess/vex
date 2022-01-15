// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 23:44:46

package vex

import (
	"runtime"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewClientPool$
func TestNewClientPool(t *testing.T) {
	server := NewServer()
	server.RegisterHandler(1, func(req []byte) (rsp []byte, err error) {
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

	pool, err := NewClientPool("tcp", "127.0.0.1:5837", 64)
	if err != nil {
		t.Fatal("new client pool failed", err)
	}
	defer pool.Close()

	for i := 0; i < 512; i++ {

		client := pool.Get()
		response, err := client.Do(1, []byte("test"))
		if err != nil {
			pool.Put(client)
			t.Fatalf("do command %d failed with %+v", i, err)
		}

		if string(response) != "test" {
			pool.Put(client)
			t.Fatalf("response %s is wrong", string(response))
		}
		pool.Put(client)
	}
}
