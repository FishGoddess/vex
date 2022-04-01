// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"errors"
	"sync"

	"github.com/FishGoddess/vex"
)

var (
	errClientPoolClosed = errors.New("vex: client pool is closed")
	errClientPoolFull   = errors.New("vex: client pool is full")
)

// State stores all states of Pool.
type State struct {
	// Opened is the opened count of connections.
	Opened int

	// Idle is the idle count of connections.
	Idle int
}

// Pool is the pool of client.
type Pool struct {
	// config stores all configuration of Pool.
	config config

	// state stores all states of Pool.
	state State

	// clients stores all unused connections.
	clients chan *poolClient

	// factory is a factory function for creating a new Client.
	factory func() (vex.Client, error)

	closed bool
	lock   sync.RWMutex
}

// NewPool returns a client pool storing some clients.
func NewPool(factory func() (vex.Client, error), opts ...Option) *Pool {
	config := newDefaultConfig()
	config.applyOptions(opts)
	return &Pool{
		config:  config,
		clients: make(chan *poolClient, config.maxOpened),
		factory: factory,
	}
}

// put stores a client to pool.
func (cp *Pool) put(client *poolClient) {
	cp.clients <- client
}

// newClient returns a new Client.
func (cp *Pool) newClient() (vex.Client, error) {
	client, err := cp.factory()
	if err != nil {
		return nil, err
	}

	cp.state.Opened++
	return wrapClient(cp, client), nil
}

// Get returns a client for use.
func (cp *Pool) Get() (vex.Client, error) {
	cp.lock.RLock()
	defer cp.lock.RUnlock()

	if cp.closed {
		return nil, errClientPoolClosed
	}

	select {
	case client := <-cp.clients:
		return client, nil // Try to get an idle client and succeed.
	default:
		// Pool is full, returns an error.
		if cp.config.fullStrategy.Failed() {
			return nil, errClientPoolFull
		}

		// Pool is full, returns a new client.
		if cp.config.fullStrategy.New() {
			return cp.newClient()
		}

		// Pool is full, blocks util pool has an idle client.
		return <-cp.clients, nil
	}
}

// State returns all states of client pool.
func (cp *Pool) State() State {
	cp.lock.RLock()
	defer cp.lock.RUnlock()
	return cp.state
}

// Close closes pool and releases all resources.
func (cp *Pool) Close() error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if cp.closed {
		return nil
	}

	for i := 0; i < cp.config.maxOpened; i++ {
		client, ok := <-cp.clients
		if !ok {
			continue
		}

		if err := client.client.Close(); err != nil {
			return err
		}
	}

	cp.closed = true
	close(cp.clients)
	return nil
}
