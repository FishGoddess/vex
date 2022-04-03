// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

const (
	// address is the address of server.
	address = "127.0.0.1:5837"

	// benchmarkPacketType is the packet type of benchmark.
	benchmarkPacketType = 1

	// benchmarkRequestBody is the request body of benchmark.
	benchmarkRequestBody = "我是水不要鱼，希望大家可以支持开源，支持国产，不管目前有啥问题，都可以用一种理性的长远目光看待~"

	// loop is the loop of test.
	loop = 100000
)

func newClient() vex.Client {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		panic(err)
	}
	return client
}

func newClientPool(maxConnected uint64) *pool.Pool {
	return pool.NewPool(func() (vex.Client, error) {
		return vex.NewClient("tcp", address)
	}, vex.WithMaxConnected(maxConnected), vex.WithBlockOnLimit())
}

func newServer() *vex.Server {
	server := vex.NewServer()
	server.RegisterPacketHandler(benchmarkPacketType, func(requestBody []byte) (responseBody []byte, err error) {
		return requestBody, nil
	})

	go func() {
		err := server.ListenAndServe("tcp", address)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	return server
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
// BenchmarkServer-16        156993              7517 ns/op             352 B/op          6 allocs/op
func BenchmarkServer(b *testing.B) {
	server := newServer()
	defer server.Close()

	client := newClient()
	defer client.Close()

	body := []byte(benchmarkRequestBody)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Send(benchmarkPacketType, body)
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
	body := []byte(benchmarkRequestBody)
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)

		func() {
			defer wg.Done()

			body, err := client.Send(benchmarkPacketType, body)
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

	clientPool := newClientPool(16)
	defer clientPool.Close()

	//go func() {
	//	for {
	//		fmt.Printf("%+v\n", clientPool.State())
	//		time.Sleep(10 * time.Millisecond)
	//	}
	//}()

	var wg sync.WaitGroup
	body := []byte(benchmarkRequestBody)
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			client, err := clientPool.Get()
			if err != nil {
				t.Error(err)
				return
			}
			defer client.Close()

			body, err := client.Send(benchmarkPacketType, body)
			if err != nil {
				t.Error(err, body)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("Taken time is %s!\n", time.Since(beginTime).String())
	t.Logf("%+v\n", clientPool.State())
}
