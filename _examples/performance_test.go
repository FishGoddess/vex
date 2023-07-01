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

type benchmarkHandler struct {
	read  bool
	write bool
}

func newBenchmarkHandler(read bool, write bool) vex.Handler {
	return &benchmarkHandler{
		read:  read,
		write: write,
	}
}

func (bh *benchmarkHandler) Handle(ctx context.Context, reader io.Reader, writer io.Writer) {
	var wg sync.WaitGroup

	if bh.read {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, len(benchmarkPacket))
			for {
				_, err := reader.Read(buf)
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Error(err, "server read")
				}
			}
		}()
	}

	if bh.write {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				_, err := writer.Write(benchmarkPacket)
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Error(err, "server write")
				}
			}
		}()
	}

	wg.Wait()
}

func newBenchmarkClient(address string) vex.Client {
	client, err := vex.NewClient(address)
	if err != nil {
		panic(err)
	}

	return client
}

func newBenchmarkServer(address string, read bool, write bool) vex.Server {
	server := vex.NewServer(address, newBenchmarkHandler(read, write), vex.WithCloseTimeout(10*time.Second))

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	return server
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkClientReadServerWrite$ -benchtime=1s
func BenchmarkClientReadServerWrite(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address, true, true)
	defer server.Close()

	client := newBenchmarkClient(address)
	defer func() {
		log.Info("client close")
		client.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	buf := make([]byte, len(benchmarkPacket))
	for i := 0; i < b.N; i++ {
		_, err := client.Read(buf)
		if err != nil {
			b.Error(i, err)
		}
	}
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkClientWriteServerRead$ -benchtime=1s
// BenchmarkClientWriteServerRead-12         287038              4232 ns/op               0 B/op          0 allocs/op
func BenchmarkClientWriteServerRead(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address, true, false)
	defer server.Close()

	client := newBenchmarkClient(address)
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.Write(benchmarkPacket)
		if err != nil {
			b.Error(i, err)
		}
	}
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkClientServerExchange$ -benchtime=1s
func BenchmarkClientServerExchange(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address, true, true)
	defer server.Close()

	client := newBenchmarkClient(address)
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()

	buf := make([]byte, len(benchmarkPacket))
	for i := 0; i < b.N; i++ {
		_, err := client.Read(buf)
		if err != nil {
			b.Error(i, err)
		}

		_, err = client.Write(benchmarkPacket)
		if err != nil {
			b.Error(i, err)
		}
	}
}

func calculateRPS(loop int, cost time.Duration) float64 {
	return math.Round(float64(loop) * float64(time.Second) / float64(cost))
}

// go test ./_examples/performance_test.go -v -run=^TestClientRPS$
func TestClientRPS(t *testing.T) {
	//addresses := []string{"127.0.0.1:6789", "127.0.0.1:7890", "127.0.0.1:8901", "127.0.0.1:9012"}
	addresses := []string{"127.0.0.1:6789"}

	servers := make([]vex.Server, 0, len(addresses))
	for _, address := range addresses {
		servers = append(servers, newBenchmarkServer(address, true, false))
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

	poolSize := uint64(1)
	clientPool := pool.New(dial, pool.WithConnections(poolSize))
	defer clientPool.Close()

	go func() {
		for {
			t.Logf("%+v", clientPool.Status())
			time.Sleep(100 * time.Millisecond)
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
				return
			}
		}()
	}

	wg.Wait()
	cost := time.Since(beginTime)

	t.Logf("%+v", clientPool.Status())
	t.Logf("PoolSize is %d, took %s, rps is %.0f!\n", poolSize, cost, calculateRPS(loop, cost))
}

// go test ./_examples/performance_test.go -v -run=^TestClientPoolRPS$
func TestClientPoolRPS(t *testing.T) {
	//addresses := []string{"127.0.0.1:6789", "127.0.0.1:7890", "127.0.0.1:8901", "127.0.0.1:9012"}
	addresses := []string{"127.0.0.1:6789"}

	servers := make([]vex.Server, 0, len(addresses))
	for _, address := range addresses {
		servers = append(servers, newBenchmarkServer(address, true, false))
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

	go func() {
		for {
			t.Logf("%+v", clientPool.Status())
			time.Sleep(100 * time.Millisecond)
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
				return
			}
		}()
	}

	wg.Wait()
	cost := time.Since(beginTime)

	t.Logf("%+v", clientPool.Status())
	t.Logf("PoolSize is %d, took %s, rps is %.0f!\n", poolSize, cost, calculateRPS(loop, cost))
}
