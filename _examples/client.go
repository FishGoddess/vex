// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	msg := []byte("hello")
	if _, err := client.Write(msg); err != nil {
		panic(err)
	}

	var buf [1024]byte
	n, err := client.Read(buf[:])
	if err != nil {
		panic(err)
	}

	fmt.Println("Received:", string(buf[:n]))
}
