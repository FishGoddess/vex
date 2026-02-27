// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"sync"
)

var contextPool = sync.Pool{
	New: func() any {
		return new(Context)
	},
}

type Context struct {
	context.Context
}

func acquireContext(parentCtx context.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = parentCtx
	return ctx
}

func releaseContext(ctx *Context) {
	ctx.Context = nil
	contextPool.Put(ctx)
}
