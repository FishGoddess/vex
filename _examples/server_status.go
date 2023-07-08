// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"time"

	"github.com/FishGoddess/vex"
)

func statusHandle(ctx *vex.Context) {
	var buf [1024]byte
	for {
		n, err := ctx.Read(buf[:])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println("Received:", string(buf[:n]))

		if _, err = ctx.Write(buf[:n]); err != nil {
			panic(err)
		}

		// Do some expensive things.
		time.Sleep(2 * time.Second)
	}
}

func watchStatus(server vex.Server) {
	for {
		fmt.Printf("%+v\n", server.Status())
		time.Sleep(time.Second)
	}
}

func main() {
	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	// By default, we set this value to 4096 which may be universal.
	// Use WithMaxConnections to limit the connections connected from clients.
	server := vex.NewServer("127.0.0.1:6789", statusHandle, vex.WithName("status"), vex.WithMaxConnections(1))
	defer server.Close()

	// Watching the status of server.
	go watchStatus(server)

	// Use Serve() to begin serving.
	// Press ctrl+c/command+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
