// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

const (
	// address is the address of server.
	address = "127.0.0.1:5837"

	// benchmarkTag is the command of benchmark.
	benchmarkTag = byte(1)

	// loop is the loop of test.
	loop = 100000
)

func newServer() *vex.Server {
	server := vex.NewServer()
	server.RegisterHandler(benchmarkTag, func(req []byte) (rsp []byte, err error) {
		return req, nil
	})

	go func() {
		err := server.ListenAndServe("tcp", address)
		if err != nil {
			panic(err)
		}
	}()

	return server
}

func newClient() vex.Client {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		panic(err)
	}
	return client
}

func newClientPool(maxConnections int) *vex.ClientPool {
	pool, err := vex.NewClientPool(maxConnections, func() (vex.Client, error) {
		return vex.NewClient("tcp", address)
	})
	if err != nil {
		panic(err)
	}
	return pool
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
// BenchmarkServer-16        187090              6632 ns/op              32 B/op          6 allocs/op
func BenchmarkServer(b *testing.B) {
	server := newServer()
	defer server.Close()

	client := newClient()
	defer client.Close()

	req := []byte("req")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Do(benchmarkTag, req)
		if err != nil {
			b.Error(err)
		}
	}
}

// go test ./_examples/performance_test.go -v -run=^TestRPS$
func TestRPS(t *testing.T) {
	server := newServer()
	defer server.Close()

	client := newClient()
	defer client.Close()

	var wg sync.WaitGroup
	req := []byte("req")
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)

		func() {
			defer wg.Done()

			body, err := client.Do(benchmarkTag, req)
			if err != nil {
				t.Error(err, body)
			}
		}()
	}

	wg.Wait()
	t.Logf("Taken time is %s!\n", time.Since(beginTime).String())
}

// go test ./_examples/performance_test.go -v -run=^TestRPSWithPool$
func TestRPSWithPool(t *testing.T) {
	server := newServer()
	defer server.Close()

	pool := newClientPool(64)
	defer pool.Close()

	var wg sync.WaitGroup
	req := []byte("req")
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			client := pool.Get()
			defer client.Close()

			body, err := client.Do(benchmarkTag, req)
			if err != nil {
				t.Error(err, body)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("Taken time is %s!\n", time.Since(beginTime).String())
}
