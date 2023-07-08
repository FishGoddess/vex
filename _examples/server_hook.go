// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/FishGoddess/vex"
)

func hookHandle(ctx *vex.Context) {
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
	}
}

func beforeServing(address string) {
	fmt.Println("before serving", address)
}

func afterServing(address string) {
	fmt.Println("after serving", address)
}

func beforeHandling(ctx *vex.Context) {
	done := false
	select {
	case <-ctx.Done():
		done = true
	default:
	}

	fmt.Printf("before handling local %s remote %s done %+v\n", ctx.LocalAddr(), ctx.RemoteAddr(), done)
}

func afterHandling(ctx *vex.Context) {
	done := false
	select {
	case <-ctx.Done():
		done = true
	default:
	}

	fmt.Printf("after handling local %s remote %s done %+v\n", ctx.LocalAddr(), ctx.RemoteAddr(), done)
}

func beforeClosing(address string) {
	fmt.Println("before closing", address)
}

func afterClosing(address string) {
	fmt.Println("after closing", address)
}

func main() {
	opts := []vex.Option{
		// We can give our server a name like "hook" so we can see it in logs.
		vex.WithName("hook"),

		// These are all server hooks we provided.
		// See the comments of functions then you will find out what they mean.
		vex.WithBeforeServing(beforeServing),
		vex.WithAfterServing(afterServing),
		vex.WithBeforeHandling(beforeHandling),
		vex.WithAfterHandling(afterHandling),
		vex.WithBeforeClosing(beforeClosing),
		vex.WithAfterClosing(afterClosing),
	}

	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	server := vex.NewServer("127.0.0.1:6789", hookHandle, opts...)

	// Use Serve() to begin serving.
	// Press ctrl+c/command+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
