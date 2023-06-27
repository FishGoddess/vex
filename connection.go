// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"net"
	"time"
)

func setupConn(conf *Config, conn *net.TCPConn) error {
	now := time.Now()
	readDeadline := now.Add(conf.ReadTimeout)
	writeDeadline := now.Add(conf.WriteTimeout)

	if err := conn.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(writeDeadline); err != nil {
		return err
	}

	if err := conn.SetReadBuffer(conf.ReadBufferSize); err != nil {
		return err
	}

	if err := conn.SetWriteBuffer(conf.WriteBufferSize); err != nil {
		return err
	}

	return nil
}
