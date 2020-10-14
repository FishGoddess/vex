// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/14 00:23:27

package vex

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewServer$
func TestNewServer(t *testing.T) {

	server := NewServer()
	server.RegisterHandler("test", func(ctx *Context) {

		if len(ctx.Args()) != 2 {
			t.Fatalf("Args length should be 2, but found %d!", len(ctx.Args()))
		}

		one := string(ctx.Arg(0))
		two := string(ctx.Arg(1))
		if one != "123" || two != "456" {
			t.Fatalf("The arg is incorrect! They are %s, %s!", one, two)
		}

		t.Log(one, two)

		_, err := ctx.Write([]byte("Test!"))
		if err != nil {
			t.Fatal(err)
		}
	})
	defer server.Close()

	go func() {
		err := server.ListenAndServe("tcp", ":5837")
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Second)

	client, err := NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	response, err := client.Do("test", [][]byte{
		[]byte("123"),
		[]byte("456"),
	})

	if err != nil {
		t.Fatal(err)
	}

	resp := string(response)
	if resp != "Test!" {
		t.Fatalf("The response %s is incorrect!", resp)
	}

	t.Log(resp)
}
