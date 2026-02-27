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
	client, err := vex.NewClient("127.0.0.1:9876")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx := context.Background()
	data := []byte("落得湖面圆满月，独守湖边酒哀愁")

	received, err := client.Send(ctx, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("client send: %s\n", data)
	fmt.Printf("client receive: %s\n", received)
}
