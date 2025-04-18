// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"testing"
)

// go test -v -cover -run=^TestWithFastFailed$
func TestWithFastFailed(t *testing.T) {
	conf := &config{fastFailed: false}
	WithFastFailed()(conf)

	if !conf.fastFailed {
		t.Errorf("conf.fastFailed %+v is wrong", conf.fastFailed)
	}
}
