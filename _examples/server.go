// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/FishGoddess/vex"
)

type handler struct{}

func (handler) Handle(conn *vex.Connection) {
	defer conn.Flush()

	buf := make([]byte, 0, 1024)
	for {
		n, err := conn.Read(buf)
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
	server := vex.NewServer("127.0.0.1:6789")
	server.Handle(handler{})

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
