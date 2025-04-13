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
	ErrPoolIsFull   = errors.New("vex: pool is full")
	ErrPoolIsClosed = errors.New("vex: pool is closed")
)

// DialFunc is a function dials to somewhere and returns a client.
// Returns an error if failed.
type DialFunc func() (vex.Client, error)

type Status struct {
	// Limit is the limit of connected clients.
	Limit uint64 `json:"limit"`

	// Connected is the count of connected clients.
	Connected uint64 `json:"connected"`

	// Idle is the count of idle clients.
	Idle uint64 `json:"idle"`

	// Waiting is the count of waiting for a client.
	Waiting uint64 `json:"waiting"`
}

type Pool struct {
	clients *rego.Pool[vex.Client]
}

func New(limit uint64, dial DialFunc, opts ...Option) *Pool {
	regoOpts := newRegoOptions(opts)

	acquire := func() (vex.Client, error) {
		return dial()
	}

	release := func(client vex.Client) error {
		return client.Close()
	}

	pool := &Pool{
		clients: rego.New(limit, acquire, release, regoOpts...),
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

	regoOpts = append(regoOpts, rego.WithPoolFullErr(func(ctx context.Context) error {
		return ErrPoolIsFull
	}))

	regoOpts = append(regoOpts, rego.WithPoolClosedErr(func(ctx context.Context) error {
		return ErrPoolIsClosed
	}))

	return regoOpts
}

func (p *Pool) Put(client vex.Client) error {
	return p.clients.Put(client)
}

func (p *Pool) Take(ctx context.Context) (vex.Client, error) {
	return p.clients.Take(ctx)
}

// Status returns the status of the pool.
func (p *Pool) Status() Status {
	status := p.clients.Status()

	return Status{
		Limit:     status.Limit,
		Connected: status.Acquired,
		Idle:      status.Idle,
		Waiting:   status.Waiting,
	}
}

// Close closes the pool.
func (p *Pool) Close() error {
	return p.clients.Close()
}
