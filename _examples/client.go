// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FishGoddess/vex"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:9876")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx := context.Background()
	for i := range 10 {
		data := []byte(strconv.Itoa(i))
		fmt.Printf("client send: %s\n", data)

		data, err = client.Send(ctx, data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("client receive: %s\n", data)
		time.Sleep(100 * time.Millisecond)
	}
}
