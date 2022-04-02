// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

func main() {
	clientPool := pool.NewPool(func() (vex.Client, error) {
		return vex.NewClient("tcp", "127.0.0.1:5837")
	})

	for i := 0; i < 10; i++ {
		client, err := clientPool.Get()
		if err != nil {
			panic(err)
		}

		responseBody, err := client.Send(1, []byte("client pool test"))
		if err != nil {
			client.Close()
			panic(err)
		}

		client.Close()
		fmt.Println(string(responseBody))
		fmt.Printf("%+v\n", clientPool.State())
		time.Sleep(time.Second)
	}

	fmt.Printf("%+v\n", clientPool.State())
}
