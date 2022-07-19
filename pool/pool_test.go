// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

// go test -v -cover -run=^TestNewPool$
func TestNewPool(t *testing.T) {
	server := vex.NewServer()
	server.RegisterPacketHandler(1, func(ctx context.Context, req []byte) (rsp []byte, err error) {
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

	pool := NewPool(func() (vex.Client, error) { return vex.NewClient("tcp", "127.0.0.1:5837") }, vex.WithMaxConnected(64))
	defer pool.Close()

	for i := 0; i < 512; i++ {
		client, err := pool.Get()
		if err != nil {
			t.Fatalf("get client failed with %+v", err)
		}

		response, err := client.Send(1, []byte("test"))
		if err != nil {
			client.Close()
			t.Fatalf("send packet %d failed with %+v", i, err)
		}

		if string(response) != "test" {
			client.Close()
			t.Fatalf("response %s is wrong", string(response))
		}

		client.Close()
	}
}
