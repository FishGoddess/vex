// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

type config struct {
	limit      uint64
	fastFailed bool
}

func newConfig() config {
	return config{
		limit:      128,
		fastFailed: false,
	}
}

type Option func(conf *config)

func (o Option) ApplyTo(conf *config) {
	o(conf)
}

// WithLimit sets limit to config.
func WithLimit(limit uint64) Option {
	return func(conf *config) {
		conf.limit = limit
	}
}

// WithFastFailed sets fastFailed to config.
func WithFastFailed() Option {
	return func(conf *config) {
		conf.fastFailed = true
	}
}
