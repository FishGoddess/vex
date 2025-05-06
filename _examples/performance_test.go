// Copyright 2025 FishGoddess. All rights reserved.
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
	//benchmarkData = []byte("我是水不要鱼，希望大家可以支持开源，支持国产，不管目前有啥问题，都可以用一种理性的长远目光看待~")
	benchmarkData = make([]byte, 1024)
)

func benchmarkHandle(ctx *vex.Context) {
	buf := make([]byte, len(benchmarkData))
	for {
		_, err := ctx.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err, "server read failed")
		}

		_, err = ctx.Write(benchmarkData)
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
func BenchmarkReadWrite(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address)
	defer server.Close()

	dial := func(ctx context.Context) (vex.Client, error) {
		return vex.NewClient(address)
	}

	clientPool := pool.New(1, dial)
	defer clientPool.Close()

	b.ReportAllocs()
	b.ResetTimer()

	ctx := context.Background()
	buf := make([]byte, len(benchmarkData))
	for i := 0; i < b.N; i++ {
		client, err := clientPool.Take(ctx)
		if err != nil {
			b.Fatal(err)
		}

		_, err = client.Write(benchmarkData)
		if err != nil {
			b.Fatal(err)
		}

		_, err = client.Read(buf)
		if err != nil {
			b.Fatal(err)
		}

		if err = clientPool.Put(ctx, client); err != nil {
			b.Fatal(err)
		}
	}
}

func calculateRPS(loop int, cost time.Duration) float64 {
	return math.Round(float64(loop) * float64(time.Second) / float64(cost))
}

// go test ./_examples/performance_test.go -v -run=^TestRPS$
func TestRPS(t *testing.T) {
	//addresses := []string{"127.0.0.1:6789", "127.0.0.1:7890", "127.0.0.1:8901", "127.0.0.1:9012"}
	addresses := []string{"127.0.0.1:9876"}

	servers := make([]vex.Server, 0, len(addresses))
	for _, address := range addresses {
		servers = append(servers, newBenchmarkServer(address))
	}

	defer func() {
		for _, server := range servers {
			if err := server.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}()

	index := uint64(0)
	dial := func(ctx context.Context) (vex.Client, error) {
		next := atomic.AddUint64(&index, 1)
		i := int(next) % len(addresses)
		return vex.NewClient(addresses[i])
	}

	poolSize := uint64(2)
	clientPool := pool.New(poolSize, dial)
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
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < loop; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := clientPool.Take(ctx)
			if err != nil {
				t.Fatal(err)
				return
			}

			defer clientPool.Put(ctx, client)

			_, err = client.Write(benchmarkData)
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, len(benchmarkData))
			_, err = client.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	wg.Wait()
	cost := time.Since(beginTime)

	t.Logf("%+v", clientPool.Status())
	t.Logf("PoolSize is %d, took %s, rps is %.0f!\n", poolSize, cost, calculateRPS(loop, cost))
}
