// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
)

var (
	connContextKey = struct{}{}
)

// connContext is a context with conn information.
type connContext struct {
	context.Context
	localAddr  net.Addr
	remoteAddr net.Addr
}

// wrapContext wraps ctx with conn.
func wrapContext(ctx context.Context, localAddr net.Addr, remoteAddr net.Addr) context.Context {
	ctx = &connContext{
		Context:    ctx,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
	return context.WithValue(ctx, connContextKey, ctx)
}

// LocalAddr returns the local network address.
func LocalAddr(ctx context.Context) (net.Addr, bool) {
	connCtx, ok := ctx.Value(connContextKey).(*connContext)
	if !ok {
		return nil, false
	}
	return connCtx.localAddr, true
}

// RemoteAddr returns the remote network address.
func RemoteAddr(ctx context.Context) (net.Addr, bool) {
	connCtx, ok := ctx.Value(connContextKey).(*connContext)
	if !ok {
		return nil, false
	}
	return connCtx.remoteAddr, true
}
