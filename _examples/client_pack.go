// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pack"
)

func useClientWithPack(client vex.Client) {
	for i := 0; i < 10; i++ {
		msg := strconv.Itoa(i)
		if i&1 == 0 {
			msg = "error"
		}

		packet, err := pack.Send(client, 1, []byte(msg))
		fmt.Println(string(packet), err)
	}
}

func main() {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()
	useClientWithPack(client)
}
