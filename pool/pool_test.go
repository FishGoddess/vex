// Copyright 2025 FishGoddess. All rights reserved.
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
	ctx := context.Background()
	address := "127.0.0.1:10000"

	server := vex.NewServer(address, handle)
	defer server.Close()

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	runtime.Gosched()
	time.Sleep(time.Second)

	dial := func(ctx context.Context) (vex.Client, error) {
		return vex.NewClient(address)
	}

	pool := New(16, dial)
	defer pool.Close(ctx)

	data := []byte("test")
	test := func(i int) {
		client, err := pool.Take(ctx)
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

		pool.Put(ctx, client)
	}

	for i := 0; i < 1024; i++ {
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
	for i := 0; i < 1024; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			test(i)

			status := pool.Status()
			if status.Connected > status.Limit {
				t.Errorf("status.Connected %d is wrong", status.Connected)
			}

			if status.Idle > status.Limit {
				t.Errorf("status.Idle %d is wrong", status.Idle)
			}
		}(i)
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())
}
