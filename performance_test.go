// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/14 22:19:15

package vex

import "testing"

// BenchmarkServer-8           8568            139309 ns/op            8593 B/op         22 allocs/op

// go test -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
func BenchmarkServer(b *testing.B) {

	server := NewServer()
	server.RegisterHandler("test", func(ctx *Context) {
		_, err := ctx.Write([]byte("Test!"))
		if err != nil {
			b.Fatal(err)
		}
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			b.Fatal(err)
		}
	}()

	client, err := NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Do("test", [][]byte{
			[]byte("123"),
			[]byte("456"),
		})

		if err != nil {
			b.Fatal(err)
		}
	}
}
