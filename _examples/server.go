// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	server := vex.NewServer("tcp", "127.0.0.1:5837", vex.WithName("example"))
	server.RegisterPacketHandler(1, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		addr, ok := vex.RemoteAddr(ctx)
		if !ok {
			fmt.Println(string(requestBody))
		} else {
			fmt.Println(string(requestBody), "from", addr)
		}
		return []byte("server test"), nil
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
