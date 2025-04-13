// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pack"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// Use Send method to send a packet to server and receive a packet from server.
	// Try to change 'hello' to 'error' and see what happens.
	packet, err := pack.Send(client, 1, []byte("hello"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(packet))
}
