// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2022/01/15 22:04:59

package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	server := vex.NewServer()
	server.RegisterHandler(1, func(req []byte) (rsp []byte, err error) {
		fmt.Println(string(req))
		return []byte("server test"), nil
	})

	err := server.ListenAndServe("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
}
