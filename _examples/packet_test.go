// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

var (
	benchmarkData = make([]byte, 1024)
)

type BenchmarkHandler struct{}

func (BenchmarkHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	return data, nil
}

func newBenchmarkClient(address string) vex.Client {
	client, err := vex.NewClient(address)
	if err != nil {
		panic(err)
	}

	return client
}

func newBenchmarkClientPool(addresses []string) *vex.ClientPool {
	var index atomic.Int64

	dial := func(ctx context.Context) (vex.Client, error) {
		i := index.Add(1) % int64(len(addresses))
		return vex.NewClient(addresses[i])
	}

	limit := uint64(len(addresses)) * 2
	clientPool := vex.NewClientPool(limit, dial)
	return clientPool
}

func newBenchmarkServer(address string) vex.Server {
	server := vex.NewServer(address, BenchmarkHandler{})

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	return server
}

func newBenchmarkServers(addresses []string) []vex.Server {
	servers := make([]vex.Server, 0, len(addresses))
	for _, address := range addresses {
		server := newBenchmarkServer(address)
		servers = append(servers, server)
	}

	return servers
}

// go test -v -run=none -bench=^BenchmarkPacket$ -benchmem -benchtime=1s ./_examples/packet_test.go
func BenchmarkPacket(b *testing.B) {
	address := "127.0.0.1:6789"

	server := newBenchmarkServer(address)
	defer server.Close()

	client := newBenchmarkClient(address)
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()

	ctx := context.Background()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := client.Send(ctx, benchmarkData)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

// go test -v -run=none -bench=^BenchmarkPacketPool$ -benchmem -benchtime=1s ./_examples/packet_test.go
func BenchmarkPacketPool(b *testing.B) {
	addresses := []string{"127.0.0.1:6789"}

	servers := newBenchmarkServers(addresses)
	for i := range servers {
		defer servers[i].Close()
	}

	clientPool := newBenchmarkClientPool(addresses)
	defer clientPool.Close()

	b.ReportAllocs()
	b.ResetTimer()

	ctx := context.Background()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := clientPool.Send(ctx, benchmarkData)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
