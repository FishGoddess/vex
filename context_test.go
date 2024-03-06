// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"net"
	"testing"
)

// go test -v -cover -count=1 -test.cpu=1 -run=^TestSetupConn$
func TestSetupConn(t *testing.T) {
	conf := newServerConfig("127.0.0.1:6789")

	listener, err := net.Listen(network, conf.address)
	if err != nil {
		t.Error(err)
	}

	defer listener.Close()

	go func() {
		listener.Accept()
	}()

	resolved, err := net.ResolveTCPAddr(network, conf.address)
	if err != nil {
		t.Error(err)
	}

	conn, err := net.DialTCP(network, nil, resolved)
	if err != nil {
		t.Error(err)
	}

	if err = setupConn(conf, conn); err != nil {
		t.Error(err)
	}
}
