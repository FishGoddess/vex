// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"testing"
)

// go test -v -cover -run=^TestContext$
func TestContext(t *testing.T) {
	parentCtx := context.Background()

	ctx := acquireContext(parentCtx)
	if ctx.Context != parentCtx {
		t.Fatalf("got %+v != want %+v", ctx.Context, parentCtx)
	}

	releaseContext(ctx)
	if ctx.Context != nil {
		t.Fatalf("got %+v != nil", ctx.Context)
	}
}
