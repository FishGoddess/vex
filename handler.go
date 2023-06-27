// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"io"
)

type Handler interface {
	Handle(ctx *Context, reader io.Reader, writer io.Writer)
}
