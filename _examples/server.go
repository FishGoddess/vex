// Copyright 2021 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 22:39:50

package main

import "github.com/FishGoddess/vex"

func main() {

	server := vex.NewServer()
	server.RegisterHandler(1, func(args [][]byte) (body []byte, err error) {
		return []byte("test"), nil
	})

	err := server.ListenAndServe("tcp", ":5837")
	if err != nil {
		panic(err)
	}
}
