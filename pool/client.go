// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "github.com/FishGoddess/vex"

type poolClient struct {
	pool   *Pool
	client vex.Client
}

func newPoolClient(pool *Pool, client vex.Client) vex.Client {
	return &poolClient{
		pool:   pool,
		client: client,
	}
}

func (pc *poolClient) closeUnderlying() error {
	return pc.client.Close()
}

// Read reads data to p.
// See io.Reader.
func (pc *poolClient) Read(p []byte) (n int, err error) {
	return pc.client.Read(p)
}

// Write writes p to data.
// See io.Writer.
func (pc *poolClient) Write(p []byte) (n int, err error) {
	return pc.client.Write(p)
}

// Close closes the client and returns an error if failed.
func (pc *poolClient) Close() error {
	return pc.pool.put(pc)
}
