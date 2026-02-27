// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
	"sync"
)

var contextPool = sync.Pool{
	New: func() any {
		return new(Context)
	},
}

func acquireContext(parentCtx context.Context, conn net.Conn) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = parentCtx
	ctx.clientAddr = conn.RemoteAddr().String()
	return ctx
}

func releaseContext(ctx *Context) {
	ctx.Context = nil
	ctx.clientAddr = ""

	contextPool.Put(ctx)
}

// Context wraps a context inside and carries some attributes for using.
type Context struct {
	context.Context

	clientAddr string
}

// ClientAddr returns the client address of connection.
func (c *Context) ClientAddr() string {
	return c.clientAddr
}
