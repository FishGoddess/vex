// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"

	"github.com/FishGoddess/vex"
)

func main() {
	handle := func(ctx *vex.Context) {
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

			reply := strconv.FormatUint(rand.Uint64(), 10)
			if _, err = ctx.Write([]byte(reply)); err != nil {
				panic(err)
			}
		}
	}

	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	// Also, we can give it a name like "example" so we can see it in logs.
	server := vex.NewServer("127.0.0.1:6789", handle, vex.WithName("example"))
	defer server.Close()

	// Use Serve() to begin serving.
	// Press ctrl+c/command+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
