// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/08/02 23:32:44

package vex

// poolClient wraps client to a pool client.
type poolClient struct {
	// pool is the owner of this client.
	pool *ClientPool

	// client is the wrapped target client.
	client Client
}

func (pc *poolClient) Do(tag Tag, req []byte) (rsp []byte, err error) {
	return pc.client.Do(tag, req)
}

func (pc *poolClient) Close() error {
	pc.pool.put(pc) // TODO Double Close will cause concurrent problem.
	return nil
}

// ClientPool is the pool of client.
type ClientPool struct {
	// maxConnections is the max count of connections.
	maxConnections int

	// clients stores all unused connections.
	clients chan Client

	// newClient returns a new Client.
	newClient func() (Client, error)
}

// NewClientPool returns a client pool storing some clients.
func NewClientPool(maxConnections int, newClient func() (Client, error)) (*ClientPool, error) {
	pool := &ClientPool{
		maxConnections: maxConnections,
		clients:        make(chan Client, maxConnections),
		newClient:      newClient,
	}

	for i := 0; i < maxConnections; i++ {
		client, err := newClient()
		if err != nil {
			return nil, err
		}

		pool.put(pool.wrapClient(client))
	}

	return pool, nil
}

// wrapClient wraps client to a pool client.
func (cp *ClientPool) wrapClient(client Client) Client {
	return &poolClient{
		pool:   cp,
		client: client,
	}
}

// put stores a client to pool.
func (cp *ClientPool) put(client Client) {
	cp.clients <- client
}

// Get returns a client for use.
func (cp *ClientPool) Get() Client {
	return <-cp.clients
}

// Close closes pool and releases all resources.
func (cp *ClientPool) Close() error {
	for i := 0; i < cp.maxConnections; i++ {
		client, ok := (<-cp.clients).(*poolClient)
		if !ok {
			continue
		}

		if err := client.client.Close(); err != nil {
			return err
		}
	}

	close(cp.clients)
	return nil
}
