// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

// Option applies functions to config.
type Option func(c *config)

// WithMaxConnected sets maxConnected to config.
func WithMaxConnected(maxConnected uint64) Option {
	return func(c *config) {
		c.MaxConnected = maxConnected
	}
}

// WithMaxIdle sets maxIdle to config.
func WithMaxIdle(maxIdle uint64) Option {
	return func(c *config) {
		c.MaxIdle = maxIdle
	}
}

// WithNonBlockOnFull sets non block on full to config.
func WithNonBlockOnFull() Option {
	return func(c *config) {
		c.BlockOnFull = false
	}
}
