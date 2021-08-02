// Copyright 2021 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 22:52:15

/*
Package vex provides an easy way to use foundation for your net operations.

1. client

	client, err := vex.NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	response, err := client.Do(1, [][]byte{
		[]byte("123"), []byte("456"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(response))

2. server:

	server := vex.NewServer()
	server.RegisterHandler(1, func(args [][]byte) (body []byte, err error) {
		return []byte("test"), nil
	})

	err := server.ListenAndServe("tcp", ":5837")
	if err != nil {
		panic(err)
	}
*/
package vex

const (
	// Version is the version string representation of vex.
	Version = "v0.2.0-alpha"
)
