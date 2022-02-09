// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 23:32:44

package vex

// ClientPool is the pool of client.
type ClientPool struct {
	// maxConnections is the max count of connections.
	maxConnections int

	newClient func() (Client, error)

	// clients stores all unused connections.
	clients chan Client
}

// NewClientPool returns a client pool storing some clients.
func NewClientPool(maxConnections int, newClient func() (Client, error)) (*ClientPool, error) {
	clients := make(chan Client, maxConnections)
	for i := 0; i < maxConnections; i++ {
		client, err := newClient()
		if err != nil {
			return nil, err
		}
		clients <- client
	}

	return &ClientPool{
		maxConnections: maxConnections,
		clients:        clients,
	}, nil
}

// Get returns a client for use.
func (cp *ClientPool) Get() Client {
	return <-cp.clients
}

// Put stores a client to pool.
func (cp *ClientPool) Put(client Client) {
	cp.clients <- client
}

// Close closes pool and releases all resources.
func (cp *ClientPool) Close() error {
	for i := 0; i < cp.maxConnections; i++ {
		client := <-cp.clients
		if err := client.Close(); err != nil {
			return err
		}
	}

	close(cp.clients)
	return nil
}
