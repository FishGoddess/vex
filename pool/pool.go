// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"errors"

	"github.com/FishGoddess/rego"
	"github.com/FishGoddess/vex"
)

var (
	ErrPoolExhausted = errors.New("vex: pool is exhausted")
	ErrPoolClosed    = errors.New("vex: pool is closed")
)

// DialFunc is a function dials to somewhere and returns a client.
// Returns an error if failed.
type DialFunc func(ctx context.Context) (vex.Client, error)

type Status struct {
	// Limit is the limit of connected clients.
	Limit uint64 `json:"limit"`

	// Connected is the quantity of connected clients.
	Connected uint64 `json:"connected"`

	// Idle is the quantity of idle clients.
	Idle uint64 `json:"idle"`

	// Waiting is the quantity of waiting for a client.
	Waiting uint64 `json:"waiting"`
}

type Pool struct {
	clientPool *rego.Pool[vex.Client]
}

func New(limit uint64, dial DialFunc, opts ...Option) *Pool {
	regoOpts := newRegoOptions(opts)

	acquire := func(ctx context.Context) (vex.Client, error) {
		return dial(ctx)
	}

	release := func(ctx context.Context, client vex.Client) error {
		return client.Close()
	}

	pool := &Pool{
		clientPool: rego.New(limit, acquire, release, regoOpts...),
	}

	return pool
}

func newRegoOptions(opts []Option) []rego.Option {
	conf := newDefaultConfig()
	for _, opt := range opts {
		opt.ApplyTo(&conf)
	}

	var regoOpts []rego.Option
	if conf.fastFailed {
		regoOpts = append(regoOpts, rego.WithFastFailed())
	}

	regoOpts = append(regoOpts, rego.WithPoolExhaustedErr(func(ctx context.Context) error {
		return ErrPoolExhausted
	}))

	regoOpts = append(regoOpts, rego.WithPoolClosedErr(func(ctx context.Context) error {
		return ErrPoolClosed
	}))

	return regoOpts
}

func (p *Pool) Put(ctx context.Context, client vex.Client) error {
	return p.clientPool.Put(ctx, client)
}

func (p *Pool) Take(ctx context.Context) (vex.Client, error) {
	return p.clientPool.Take(ctx)
}

// Status returns the status of the pool.
func (p *Pool) Status() Status {
	status := p.clientPool.Status()

	return Status{
		Limit:     status.Limit,
		Connected: status.Acquired,
		Idle:      status.Idle,
		Waiting:   status.Waiting,
	}
}

// Close closes the pool.
func (p *Pool) Close(ctx context.Context) error {
	return p.clientPool.Close(ctx)
}
