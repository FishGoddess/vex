// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:56:35

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

// BenchmarkServer-16        112126             12759 ns/op             144 B/op         11 allocs/op

const (
	// address is the address of server.
	address = "127.0.0.1:5837"

	// benchmarkCommand is the command of benchmark.
	benchmarkCommand = byte(1)

	// loop is the loop of test.
	loop = 100000
)

func newServer() *vex.Server {
	server := vex.NewServer()

	resp := []byte("test")
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (body []byte, err error) {
		return resp, nil
	})

	go func() {
		err := server.ListenAndServe("tcp", address)
		if err != nil {
			panic(err)
		}
	}()

	return server
}

func newClient() *vex.Client {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		panic(err)
	}
	return client
}

func newClientPool(maxConnections int) *vex.ClientPool {
	pool, err := vex.NewClientPool("tcp", address, maxConnections)
	if err != nil {
		panic(err)
	}
	return pool
}

// go test ./_examples/performance_test.go -v -run=^TestServerRPS$
func TestServerRPS(t *testing.T) {
	server := newServer()
	defer server.Close()

	client := newClient()
	defer client.Close()

	param1 := []byte("one")
	param2 := []byte("two")

	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		body, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
		if err != nil {
			t.Error(err, body)
		}
	}

	t.Logf("Taken time is %s!\n", time.Since(beginTime).String())
}

// go test ./_examples/performance_test.go -v -run=^TestServerRPSWithPool$
func TestServerRPSWithPool(t *testing.T) {
	server := newServer()
	defer server.Close()

	pool := newClientPool(64)
	defer pool.Close()

	param1 := []byte("one")
	param2 := []byte("two")

	wg := &sync.WaitGroup{}
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := pool.Get()
			defer pool.Put(client)

			body, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
			if err != nil {
				t.Error(err, body)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("Taken time is %s!\n", time.Since(beginTime).String())
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
func BenchmarkServer(b *testing.B) {
	server := newServer()
	defer server.Close()

	client := newClient()
	defer client.Close()

	param1 := []byte("one")
	param2 := []byte("two")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
		if err != nil {
			b.Fatal(err)
		}
	}
}
