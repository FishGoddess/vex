// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pack"
)

func newRouter() *pack.Router {
	router := pack.NewRouter()

	// Use Register method to register your handler for some packets.
	router.Register(1, func(ctx context.Context, packetType pack.PacketType, requestPacket []byte) (responsePacket []byte, err error) {
		msg := string(requestPacket)
		fmt.Println(msg)

		if msg == "error" {
			return nil, errors.New(msg)
		} else {
			return requestPacket, nil
		}
	})

	return router
}

func main() {
	// Create a router for packets.
	router := newRouter()

	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	server := vex.NewServer("127.0.0.1:6789", router.Handle, vex.WithName("pack"))

	// Use Serve() to begin serving.
	// Press ctrl+c/command+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
