// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/log"
	"github.com/FishGoddess/vex/pool"
)

var (
	//benchmarkPacket = []byte("我是水不要鱼，希望大家可以支持开源，支持国产，不管目前有啥问题，都可以用一种理性的长远目光看待~")
	benchmarkPacket = make([]byte, 1024)
)

func benchmarkHandle(ctx *vex.Context) {
	buf := make([]byte, len(benchmarkPacket))
	for {
		_, err := ctx.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err, "server read failed")
		}

		_, err = ctx.Write(benchmarkPacket)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err, "server write failed")
		}
	}
}

func newBenchmarkServer(address string) vex.Server {
	server := vex.NewServer(address, benchmarkHandle, vex.WithCloseTimeout(10*time.Second))

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	return server
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkReadWrite$ -benchtime=1s
// BenchmarkReadWrite-16             183592              6603 ns/op               0 B/op          0 allocs/op
func BenchmarkReadWrite(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address)
	defer server.Close()

	//client := newBenchmarkClient(address)
	//defer client.Close()

	clientPool := pool.New(pool.Dial(address), pool.WithConnections(1))
	defer clientPool.Close()

	b.ReportAllocs()
	b.ResetTimer()

	buf := make([]byte, len(benchmarkPacket))
	for i := 0; i < b.N; i++ {
		client, err := clientPool.Get(context.Background())
		if err != nil {
			b.Error(err)
		}

		_, err = client.Write(benchmarkPacket)
		if err != nil {
			b.Error(err)
		}

		_, err = client.Read(buf)
		if err != nil {
			b.Error(err)
		}

		if err = client.Close(); err != nil {
			b.Error(err)
		}
	}
}

func calculateRPS(loop int, cost time.Duration) float64 {
	return math.Round(float64(loop) * float64(time.Second) / float64(cost))
}

// go test ./_examples/performance_test.go -v -run=^TestRPS$
// PoolSize is 1, took 1.266500745s, rps is 78958
// PoolSize is 16, took 393.082456ms, rps is 254400
func TestRPS(t *testing.T) {
	//addresses := []string{"127.0.0.1:6789", "127.0.0.1:7890", "127.0.0.1:8901", "127.0.0.1:9012"}
	addresses := []string{"127.0.0.1:6789"}

	servers := make([]vex.Server, 0, len(addresses))
	for _, address := range addresses {
		servers = append(servers, newBenchmarkServer(address))
	}

	defer func() {
		for _, server := range servers {
			if err := server.Close(); err != nil {
				t.Error(err)
			}
		}
	}()

	index := uint64(0)
	dial := func() (vex.Client, error) {
		next := atomic.AddUint64(&index, 1)
		i := int(next) % len(addresses)
		return vex.NewClient(addresses[i])
	}

	poolSize := uint64(16)
	clientPool := pool.New(dial, pool.WithConnections(poolSize))
	defer clientPool.Close()

	doneCh := make(chan struct{}, 1)
	defer func() {
		doneCh <- struct{}{}
	}()

	go func() {
		for {
			select {
			case <-doneCh:
				return
			default:
				t.Logf("%+v", clientPool.Status())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	loop := 100000
	beginTime := time.Now()

	var wg sync.WaitGroup
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

			_, err = client.Write(benchmarkPacket)
			if err != nil {
				t.Error(err)
			}

			buf := make([]byte, len(benchmarkPacket))
			_, err = client.Read(buf)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()
	cost := time.Since(beginTime)

	t.Logf("%+v", clientPool.Status())
	t.Logf("PoolSize is %d, took %s, rps is %.0f!\n", poolSize, cost, calculateRPS(loop, cost))
}
