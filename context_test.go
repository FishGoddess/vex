// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
	"testing"
)

func newTestContext() (context.Context, net.Addr, net.Addr) {
	localAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", "192.168.32.1:9000")
	if err != nil {
		panic(err)
	}

	var ctx context.Context = &connContext{
		Context:    context.Background(),
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
	return context.WithValue(ctx, connContextKey, ctx), localAddr, remoteAddr
}

// go test -v -cover =^TestLocalAddr$
func TestLocalAddr(t *testing.T) {
	ctx, localAddr, _ := newTestContext()
	addr, ok := LocalAddr(ctx)
	if !ok {
		t.Errorf("get local address from ctx %+v not ok", ctx)
	}

	if addr.String() != localAddr.String() {
		t.Errorf("addr.String() %s != localAddr.String() %s", addr.String(), localAddr.String())
	}
	t.Log(addr)
}

// go test -v -cover =^TestRemoteAddr$
func TestRemoteAddr(t *testing.T) {
	ctx, _, remoteAddr := newTestContext()
	addr, ok := RemoteAddr(ctx)
	if !ok {
		t.Errorf("get remote address from ctx %+v not ok", ctx)
	}

	if addr.String() != remoteAddr.String() {
		t.Errorf("addr.String() %s != remoteAddr.String() %s", addr.String(), remoteAddr.String())
	}
	t.Log(addr)
}
