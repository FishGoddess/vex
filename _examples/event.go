// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/FishGoddess/vex"
)

// newDemoEventListener returns an event listener for demo.
func newDemoEventListener() vex.EventListener {
	listener := vex.NewLogEventListener()
	return vex.EventListener{
		OnServerStart: func(event vex.ServerStartEvent) {
			listener.CallOnServerStart(event)
			fmt.Println("OnServerStart...")
		},
		OnServerShutdown: func(event vex.ServerShutdownEvent) {
			listener.CallOnServerShutdown(event)
			fmt.Println("OnServerShutdown...")
		},
		OnServerGotConnected: func(event vex.ServerGotConnectedEvent) {
			listener.CallOnServerGotConnected(event)
			fmt.Println("OnServerGotConnected...")
		},
		OnServerGotDisconnected: func(event vex.ServerGotDisconnectedEvent) {
			listener.CallOnServerGotDisconnected(event)
			fmt.Println("OnServerGotDisconnected...")
		},
	}
}

func main() {
	var server *vex.Server

	// We add a default event listener which will log some common events.
	//server = vex.NewServer()

	// Also, you can customize your own event listener by implementing methods in EventListener.
	server = vex.NewServer("tcp", "127.0.0.1:5837", vex.WithEventListener(newDemoEventListener()))

	// Try to connect this server and toggle some events to see what happens.
	server.RegisterPacketHandler(1, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		return requestBody, nil
	})

	// Listen and serve!
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
