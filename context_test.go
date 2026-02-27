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

	clientAddr := conn.RemoteAddr().String()
	if ctx.clientAddr != clientAddr {
		t.Fatalf("got %s != want %s", ctx.clientAddr, clientAddr)
	}

	releaseContext(ctx)
	if ctx.Context != nil {
		t.Fatalf("got %+v != nil", ctx.Context)
	}

	if ctx.clientAddr != "" {
		t.Fatalf("got %+v != ''", ctx.clientAddr)
	}
}
