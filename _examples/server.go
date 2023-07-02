// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/FishGoddess/vex"
)

type handler struct{}

func (handler) Handle(ctx context.Context, conn *vex.Connection) {
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}

		fmt.Println("Received:", string(buf[:n]))

		reply := strconv.FormatUint(rand.Uint64(), 10)
		if _, err = conn.Write([]byte(reply)); err != nil {
			panic(err)
		}
	}
}

func main() {
	// Create a server listening on 127.0.0.1:6789 and set a handler to it.
	server := vex.NewServer("127.0.0.1:6789", handler{})

	// Use Serve() to begin serving.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
