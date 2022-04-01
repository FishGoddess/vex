// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

// poolClient wraps client to a pool client.
type poolClient struct {
	// pool is the owner of this client.
	pool *ClientPool

	// client is the target client to be wrapped.
	client Client
}

// wrapClient wraps client to a pool client.
func wrapClient(pool *ClientPool, client Client) Client {
	return &poolClient{
		pool:   pool,
		client: client,
	}
}

// Send sends a packet with requestBody to server and returns responseBody responded from server.
func (pc *poolClient) Send(packetType PacketType, requestBody []byte) (responseBody []byte, err error) {
	return pc.client.Send(packetType, requestBody)
}

// Close closes current client.
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

		pool.put(wrapClient(pool, client))
	}

	return pool, nil
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
