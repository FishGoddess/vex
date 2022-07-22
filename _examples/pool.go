// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

func main() {
	clientPool := pool.NewPool(func() (vex.Client, error) {
		return vex.NewClient("tcp", "127.0.0.1:5837")
	}, pool.WithMaxConnected(64), pool.WithMaxIdle(16))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := clientPool.Get(context.Background())
			if err != nil {
				panic(err)
			}
			defer client.Close()

			responseBody, err := client.Send(1, []byte("client pool test"))
			if err != nil {
				panic(err)
			}

			fmt.Println(string(responseBody))
			fmt.Printf("%+v\n", clientPool.State())
		}()
	}

	wg.Wait()
	fmt.Printf("%+v\n", clientPool.State())
}
