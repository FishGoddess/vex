// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

func main() {
	clientPool := pool.NewPool(func() (vex.Client, error) {
		return vex.NewClient("tcp", "127.0.0.1:5837")
	})

	client, err := clientPool.Get()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	responseBody, err := client.Send(1, []byte("client test"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(responseBody))
	fmt.Println(clientPool.State())
}
