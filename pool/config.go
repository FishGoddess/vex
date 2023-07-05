// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

type Config struct {
	// maxConnected is the max-opened count of connections.
	maxConnected uint64

	// maxIdle is the max-idle count of connections.
	maxIdle uint64

	// blockOnFull means getting clients from pool will block if connected is greater than max connected.
	blockOnFull bool
}

func newDefaultConfig() *Config {
	return &Config{
		maxConnected: 256,
		maxIdle:      256,
		blockOnFull:  true,
	}
}

func (c *Config) ApplyOptions(opts []Option) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}
