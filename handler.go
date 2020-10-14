// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/12 23:38:56

package vex

type Handler func(ctx *Context)

func readRequestErrorHandler(ctx *Context) {
	ctx.WriteError("Failed to read request!")
}

func notFoundErrorHandler(ctx *Context) {
	ctx.WriteError("Failed to find this handler!")
}