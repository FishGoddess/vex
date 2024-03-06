// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pool"
)

func main() {
	// We provide a client pool for limiting the number of clients.
	// It needs a dial function to create a new client when it needs.
	// You can pass any client options to Dial to create the client as usual.
	clientPool := pool.New(pool.Dial("127.0.0.1:6789"))
	defer clientPool.Close()

	client, err := clientPool.Get(context.Background())
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// Use client as usual.
	// Also, you can customize your dial function:
	dial := func() (vex.Client, error) {
		// You can do anything you want to customize the client.
		return vex.NewClient("127.0.0.1:6789", vex.WithReadBufferSize(4096), vex.WithWriteBufferSize(4096))
	}

	pool.New(dial)

	// The pool has some default configurations.
	// If you want to change them, see pool.Option.
	pool.New(dial, pool.WithLimit(16))
}
