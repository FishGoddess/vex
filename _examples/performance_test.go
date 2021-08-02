// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 18:56:35

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

// BenchmarkServer-16        112126             12759 ns/op             144 B/op         11 allocs/op

const (
	// dataSize is the data size of test.
	dataSize = 10000

	// concurrency is the concurrency of test.
	concurrency = 10000

	// benchmarkCommand is the command of benchmark.
	benchmarkCommand = byte(1)
)

// testTask is a wrapper wraps task to testTask.
func testTask(task func(no int)) string {
	beginTime := time.Now()
	for i := 0; i < dataSize; i++ {
		task(i)
	}
	return time.Now().Sub(beginTime).String()
}

// testTaskConcurrent is a wrapper wraps task to testTask with concurrency.
func testTaskConcurrent(task func(no int)) string {

	wg := &sync.WaitGroup{}
	beginTime := time.Now()
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

// go test ./_examples/performance_test.go -v -run=^TestServerRPS$
func TestServerRPS(t *testing.T) {

	resp := []byte("test")
	param1 := []byte("one")
	param2 := []byte("two")

	server := vex.NewServer()
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (body []byte, err error) {
		return resp, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			panic(err)
		}
	}()

	client, err := vex.NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	takenTime := testTask(func(no int) {
		body, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
		if err != nil {
			t.Fatal(err, body)
		}
	})

	t.Logf("Taken time is %s!\n", takenTime)
}

// go test ./_examples/performance_test.go -v -run=^TestServerRPSWithPool$
func TestServerRPSWithPool(t *testing.T) {

	resp := []byte("test")
	param1 := []byte("one")
	param2 := []byte("two")

	server := vex.NewServer()
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (body []byte, err error) {
		return resp, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			panic(err)
		}
	}()

	pool, err := vex.NewClientPool("tcp", "127.0.0.1:5837", 64)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	takenTime := testTaskConcurrent(func(no int) {
		client := pool.Get()
		defer pool.Put(client)

		body, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
		if err != nil {
			t.Fatal(err, body)
		}
	})

	t.Logf("Taken time is %s!\n", takenTime)
}

// go test ./_examples/performance_test.go -v -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
func BenchmarkServer(b *testing.B) {

	resp := []byte("test")
	param1 := []byte("one")
	param2 := []byte("two")

	server := vex.NewServer()
	server.RegisterHandler(benchmarkCommand, func(args [][]byte) (body []byte, err error) {
		return resp, nil
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			panic(err)
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
		_, err := client.Do(benchmarkCommand, [][]byte{param1, param2})
		if err != nil {
			b.Fatal(err)
		}
	}
}
