// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:56:35

package main

import (
	"testing"

	"github.com/FishGoddess/vex"
)

// BenchmarkServer-8          45967             23805 ns/op             144 B/op         12 allocs/op

const (
	benchmarkCommand = byte(1)
)

// go test -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
func BenchmarkServer(b *testing.B) {

	server := vex.NewServer()
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (reply byte, body []byte, err error) {
		return vex.SuccessReply, []byte("test"), nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			b.Fatal(err)
		}
	}()

	client, err := vex.NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Do(benchmarkCommand, [][]byte{
			[]byte("123"),
			[]byte("456"),
		})

		if err != nil {
			b.Fatal(err)
		}
	}
}
