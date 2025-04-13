// Copyright 2025 FishGoddess. All rights reserved.
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
	dial := func() (vex.Client, error) {
		// You can do anything you want to customize the client.
		return vex.NewClient("127.0.0.1:6789", vex.WithReadBufferSize(4096), vex.WithWriteBufferSize(4096))
	}

	clientPool := pool.New(4, dial)
	defer clientPool.Close()

	client, err := clientPool.Take(context.Background())
	if err != nil {
		panic(err)
	}

	defer clientPool.Put(client)
}
