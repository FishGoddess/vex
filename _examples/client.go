// Copyright 2021 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 22:40:39

package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {

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
}
