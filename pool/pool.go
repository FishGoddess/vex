// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"errors"
	"sync"

	"github.com/FishGoddess/vex"
)

var (
	errClientPoolFull   = errors.New("vex: client pool is full")
	errClientPoolClosed = errors.New("vex: client pool is closed")
)

// DialFunc is a function dials to somewhere and returns a client and error if failed.
type DialFunc func() (vex.Client, error)

// Dial returns a function which dials to address with opts.
func Dial(address string, opts ...vex.Option) DialFunc {
	return func() (vex.Client, error) {
		return vex.NewClient(address, opts...)
	}
}

type Status struct {
	// Connected is the count of connected connections.
	Connected uint64 `json:"connected"`

	// Idle is the count of idle connections.
	Idle uint64 `json:"idle"`

	// Waiting is the count of requests waiting for a client.
	Waiting uint64 `json:"waiting"`
}

type Pool struct {
	Config

	// dial is for creating a new Client.
	dial DialFunc

	clients chan *poolClient
	status  Status
	closed  bool

	lock sync.RWMutex
}

func New(dial DialFunc, opts ...Option) *Pool {
	conf := newDefaultConfig().ApplyOptions(opts)

	return &Pool{
		Config:  *conf,
		dial:    dial,
		clients: make(chan *poolClient, conf.maxConnected),
		closed:  false,
	}
}

func (p *Pool) newClient() (vex.Client, error) {
	client, err := p.dial()
	if err != nil {
		return nil, err
	}

	client = newPoolClient(p, client)
	return client, nil
}

func (p *Pool) put(client *poolClient) error {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		return client.closeUnderlying()
	}

	// Only waiting count < idle count will close the client immediately.
	if p.status.Waiting < p.status.Idle && p.status.Idle >= p.maxIdle {
		p.status.Connected--
		p.lock.Unlock()

		return client.closeUnderlying()
	}

	defer p.lock.Unlock()

	select {
	case p.clients <- client:
		p.status.Idle++
		return nil
	default:
		return client.closeUnderlying()
	}
}

// tryToGet tries to get an idle client from pool and return false if failed.
func (p *Pool) tryToGet() (*poolClient, bool) {
	select {
	case client := <-p.clients:
		return client, true
	default:
		return nil, false
	}
}

// waitToGet waits to get an idle client from pool.
// Record: Add ctx.Done() to select will cause a performance problem...
// The select will call runtime.selectgo if there are more than one case in it, and runtime.selectgo has two steps which is slow:
//
//     sellock(scases, lockorder)
//     sg := acquireSudog()
//
// We don't know what to do yet, but we think timeout mechanism should be supported even we haven't solved it.
func (p *Pool) waitToGet(ctx context.Context) (*poolClient, error) {
	select {
	case client := <-p.clients:
		if client == nil {
			return nil, errClientPoolClosed
		}

		return client, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Get gets a client from pool and returns an error if failed.
// You should call client.Close() to put a client back to the pool.
// We recommend you use a defer after getting a client.
func (p *Pool) Get(ctx context.Context) (vex.Client, error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		return nil, errClientPoolClosed
	}

	client, ok := p.tryToGet()
	if ok {
		p.status.Idle--
		p.lock.Unlock()

		if client == nil {
			return nil, errClientPoolClosed
		}

		return client, nil
	}

	if p.status.Connected < p.maxConnected {
		p.status.Connected++
		p.lock.Unlock()

		// Increase the connected and unlock before new client may cause the pool becomes full in advance.
		// So we should decrease the connected if new client failed.
		client, err := p.newClient()
		if err != nil {
			p.lock.Lock()
			p.status.Connected--
			p.lock.Unlock()

			return nil, err
		}

		return client, nil
	}

	if p.blockOnFull {
		p.status.Waiting++
		p.lock.Unlock()

		client, err := p.waitToGet(ctx)
		if err != nil {
			p.lock.Lock()
			p.status.Waiting--
			p.lock.Unlock()

			return nil, err
		}

		p.lock.Lock()
		p.status.Idle--
		p.status.Waiting--
		p.lock.Unlock()

		return client, nil
	}

	p.lock.Unlock()
	return nil, errClientPoolFull
}

// Status returns the status of the pool.
func (p *Pool) Status() Status {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.status
}

// Close closes pool and releases all resources.
func (p *Pool) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return nil
	}

	for i := uint64(0); i < p.status.Connected; i++ {
		client := <-p.clients
		if client == nil {
			continue
		}

		if err := client.closeUnderlying(); err != nil {
			return err
		}
	}

	p.status = Status{}
	p.closed = true
	close(p.clients)

	return nil
}
