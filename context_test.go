// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
	"testing"
)

// go test -v -cover -run=^TestContext$
func TestContext(t *testing.T) {
	parentCtx := context.Background()

	conn, err := net.Dial("tcp", "www.google.com:80")
	if err != nil {
		t.Fatal(err)
	}

	ctx := acquireContext(parentCtx, conn)
	if ctx.Context != parentCtx {
		t.Fatalf("got %+v != want %+v", ctx.Context, parentCtx)
	}

	localAddress := conn.LocalAddr().String()
	if ctx.localAddress != localAddress {
		t.Fatalf("got %s != want %s", ctx.localAddress, localAddress)
	}

	remoteAddress := conn.RemoteAddr().String()
	if ctx.remoteAddress != remoteAddress {
		t.Fatalf("got %s != want %s", ctx.remoteAddress, remoteAddress)
	}

	releaseContext(ctx)
	if ctx.Context != nil {
		t.Fatalf("got %+v != nil", ctx.Context)
	}

	if ctx.localAddress != "" {
		t.Fatalf("got %+v != ''", ctx.localAddress)
	}

	if ctx.remoteAddress != "" {
		t.Fatalf("got %+v != ''", ctx.remoteAddress)
	}
}
