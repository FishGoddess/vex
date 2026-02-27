// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log/slog"

	"github.com/FishGoddess/vex"
)

type EchoHandler struct{}

func (EchoHandler) Handle(ctx *vex.Context, data []byte) ([]byte, error) {
	clientAddr := ctx.ClientAddr()
	slog.Info(fmt.Sprintf("client %s send %s\n", clientAddr, data))

	data = []byte("好！！！")
	return data, nil
}

func main() {
	server := vex.NewServer("127.0.0.1:9876", EchoHandler{})
	defer server.Close()

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
