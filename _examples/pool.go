// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	dial := func(ctx context.Context) (vex.Client, error) {
		return vex.NewClient("127.0.0.1:9876")
	}

	pool := vex.NewPool(4, dial)
	defer pool.Close()

	ctx := context.Background()
	data := []byte("落得湖面月圆满，独守湖边酒哀愁")

	client, err := pool.Get(ctx)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	received, err := client.Send(ctx, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("client send: %s\n", data)
	fmt.Printf("server send: %s\n", received)
}
