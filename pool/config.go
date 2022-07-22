// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

// config stores all configuration of Pool.
type config struct {
	// MaxConnected is the max-opened count of connections.
	MaxConnected uint64

	// MaxIdle is the max-idle count of connections.
	MaxIdle uint64

	// BlockOnFull means getting clients from pool will block if connected is greater than max connected.
	BlockOnFull bool
}

// newDefaultConfig returns a default config.
func newDefaultConfig() *config {
	return &config{
		MaxConnected: 256,
		MaxIdle:      256,
		BlockOnFull:  true,
	}
}

// ApplyOptions applies opts to config.
func (c *config) ApplyOptions(opts []Option) *config {
	for _, opt := range opts {
		opt(c)
	}
	return c
}
