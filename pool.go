// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"

	"github.com/FishGoddess/rego"
)

type DialFunc func(ctx context.Context) (Client, error)

type ClientPool struct {
	conf *config

	pool *rego.Pool[Client]
}

func NewClientPool(limit uint64, dial DialFunc, opts ...Option) *ClientPool {
	conf := newConfig().apply(opts...)

	acquireClient := func(ctx context.Context) (Client, error) {
		return dial(ctx)
	}

	releaseClient := func(ctx context.Context, client Client) error {
		return client.Close()
	}

	pool := rego.New(limit, acquireClient, releaseClient)

	clientPool := &ClientPool{conf: conf, pool: pool}
	return clientPool
}

func (cp *ClientPool) Send(ctx context.Context, data []byte) ([]byte, error) {
	client, err := cp.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cp.pool.Release(ctx, client); err != nil {
			cp.conf.logger.Error("client released failed", "err", err)
		}
	}()

	return client.Send(ctx, data)
}

func (cp *ClientPool) Close() error {
	ctx := context.Background()
	return cp.pool.Close(ctx)
}
