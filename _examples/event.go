// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/FishGoddess/vex"
)

type demoEventHandler struct {
}

func (deh *demoEventHandler) HandleEvent(ctx context.Context, e vex.Event) {
	fmt.Println("I received an event!", e)
}

func main() {
	var server *vex.Server

	// We add a default event handler which will logs serving and shutdown events.
	//server = vex.NewServer()

	// Default event handler has a name field that you can specify.
	// This name will also be logged, so you can create more than one server using the same handler and its logs can be distinguished, either.
	//server = vex.NewServer(vex.WithEventHandler(vex.NewDefaultEventHandler("mine")))

	// Also, you can customize your own event handler by implementing interface EventHandler.
	server = vex.NewServer(vex.WithEventHandler(&demoEventHandler{}))

	// Listen and serve!
	// Try to connect this server and switch the event handler to see what happens.
	server.RegisterPacketHandler(1, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		return requestBody, nil
	})

	err := server.ListenAndServe("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
}
