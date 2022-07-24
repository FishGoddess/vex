// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

const (
	// benchmarkPacketType is the packet type of benchmark.
	benchmarkPacketType = 1
)

var (
	// benchmarkRequestBody is the request body of benchmark.
	// benchmarkRequestBody = "我是水不要鱼，希望大家可以支持开源，支持国产，不管目前有啥问题，都可以用一种理性的长远目光看待~"
	benchmarkRequestBody = make([]byte, 1024)
)

func newTestClient(address string) vex.Client {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		panic(err)
	}
	return client
}

func newTestServer(address string) *vex.Server {
	server := vex.NewServer("tcp", address, vex.WithCloseTimeout(3*time.Second), vex.WithEventHandler(nil))
	server.RegisterPacketHandler(benchmarkPacketType, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		return requestBody, nil
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	return server
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
func BenchmarkServer(b *testing.B) {
	address := "127.0.0.1:5837"

	server := newTestServer(address)
	defer server.Close()

	client := newTestClient(address)
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Send(benchmarkPacketType, benchmarkRequestBody)
		if err != nil {
			b.Error(err)
		}
	}
}

func calculateRPS(loop int, taken time.Duration) float64 {
	return math.Round(float64(loop) * float64(time.Second) / float64(taken))
}

// go test ./_examples/performance_test.go -v -run=^TestServerRPS$
func TestServerRPS(t *testing.T) {
	//addresses := []string{"127.0.0.1:5837", "127.0.0.1:6837", "127.0.0.1:7837", "127.0.0.1:8837"}
	addresses := []string{"127.0.0.1:5837"}

	servers := make([]*vex.Server, 0, len(addresses))
	for _, address := range addresses {
		servers = append(servers, newTestServer(address))
	}

	defer func() {
		for _, server := range servers {
			server.Close()
		}
	}()

	index := uint64(0)
	newClient := func() (vex.Client, error) {
		i := atomic.LoadUint64(&index)
		atomic.AddUint64(&index, 1)
		return vex.NewClient("tcp", addresses[int(i)%len(addresses)])
	}

	poolSize := uint64(16)
	clientPool := pool.NewPool(newClient, pool.WithMaxConnected(poolSize), pool.WithMaxIdle(poolSize))
	defer clientPool.Close()

	//go func() {
	//	for {
	//		fmt.Printf("%+v\n", clientPool.State())
	//		time.Sleep(10 * time.Millisecond)
	//	}
	//}()

	var wg sync.WaitGroup
	loop := 100000
	beginTime := time.Now()
	for i := 0; i < loop; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			client, err := clientPool.Get(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			defer client.Close()

			body, err := client.Send(benchmarkPacketType, benchmarkRequestBody)
			if err != nil {
				t.Error(err, body)
				return
			}
		}()
	}

	wg.Wait()
	taken := time.Since(beginTime)
	t.Logf("PoolSize is %d, took %s, rps is %.0f!\n", poolSize, taken.String(), calculateRPS(loop, taken))
	t.Logf("%+v\n", clientPool.State())
}
