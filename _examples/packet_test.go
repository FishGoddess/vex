// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
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
	for i := 0; i < b.N; i++ {
		_, err := client.Send(ctx, benchmarkData)
		if err != nil {
			b.Error(err)
		}
	}
}
