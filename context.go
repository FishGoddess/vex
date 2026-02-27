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
	ctx.localAddress = conn.LocalAddr().String()
	ctx.remoteAddress = conn.RemoteAddr().String()
	return ctx
}

func releaseContext(ctx *Context) {
	ctx.Context = nil
	ctx.localAddress = ""
	ctx.remoteAddress = ""

	contextPool.Put(ctx)
}

// Context wraps a context inside and carries some attributes for using.
type Context struct {
	context.Context

	localAddress  string
	remoteAddress string
}

// LocalAddress returns the address of server.
func (c *Context) LocalAddress() string {
	return c.localAddress
}

// RemoteAddress returns the address of client.
func (c *Context) RemoteAddress() string {
	return c.remoteAddress
}
