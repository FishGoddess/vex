// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:56:35

package main

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

// BenchmarkServer-8          45967             23805 ns/op             144 B/op         12 allocs/op

const (
	concurrency = 1

	benchmarkCommand = byte(1)
)

// testTask is a wrapper wraps task to testTask.
func testTask(task func(no int)) string {

	beginTime := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(no int) {
			defer wg.Done()
			task(no)
		}(i)
	}
	wg.Wait()
	return time.Now().Sub(beginTime).String()
}

// go test -v -run=^TestVexServer$
func TestVexServer(t *testing.T) {

	server := vex.NewServer()
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (reply byte, body []byte, err error) {
		return vex.SuccessReply, []byte("test"), nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			t.Fatal(err)
		}
	}()

	client, err := vex.NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	takenTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		body, err := client.Do(benchmarkCommand, [][]byte{
			[]byte(data), []byte(data),
		})

		if err != nil {
			t.Fatal(err, body)
		}
	})

	t.Logf("Taken time is %s!\n", takenTime)
}

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
