// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Serve() error
	Close() error
}

func monitorSignals(server Server) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-signalCh
	log.Printf("vex: received signal %+v", sig)

	if err := server.Close(); err != nil {
		log.Println(err)
	}
}
