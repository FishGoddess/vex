// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

func handle(ctx *vex.Context) {
	var buf [1024]byte
	for {
		n, err := ctx.Read(buf[:])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		n, err = ctx.Write(buf[:n])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}
}

// go test -v -cover -run=^TestPool$
func TestPool(t *testing.T) {
	address := "127.0.0.1:10000"

	server := vex.NewServer(address, handle)
	defer server.Close()

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)

	pool := New(Dial(address), WithMaxConnected(16), WithMaxIdle(64))
	defer pool.Close()

	data := []byte("test")
	test := func(i int) {
		client, err := pool.Get(context.Background())
		if err != nil {
			t.Error(err)
		}

		n, err := client.Write(data)
		if err != nil {
			client.Close()
			t.Error(i, err)
		}

		if n != len(data) {
			client.Close()
			t.Errorf("n %d != len(data) %d", n, len(data))
		}

		buf := make([]byte, 64)

		n, err = client.Read(buf[:])
		if err != nil {
			client.Close()
			t.Error(err)
		}

		if n != len(data) {
			client.Close()
			t.Errorf("n %d != len(data) %d", n, len(data))
		}

		if string(buf[:n]) != string(data) {
			client.Close()
			t.Errorf("buf %s != data %s", buf[:n], data)
		}

		client.Close()
	}

	for i := 0; i < 4096; i++ {
		test(i)

		status := pool.Status()
		if status.Connected != 1 {
			t.Errorf("status.Connected %d is wrong", status.Connected)
		}

		if status.Idle != 1 {
			t.Errorf("status.Idle %d is wrong", status.Idle)
		}
	}

	t.Logf("%+v", pool.Status())

	var wg sync.WaitGroup
	for i := 0; i < 4096; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			test(i)

			status := pool.Status()
			if status.Connected > pool.maxConnected {
				t.Errorf("status.Connected %d is wrong", status.Connected)
			}

			if status.Idle > pool.maxIdle {
				t.Errorf("status.Idle %d is wrong", status.Idle)
			}
		}(i)
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())
}
