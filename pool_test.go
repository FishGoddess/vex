// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

// go test -v -cover -run=^TestNewPool$
func TestPool(t *testing.T) {
	ctx := context.Background()

	address, done, err := runTestServer()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	var clientPtr atomic.Pointer[Client]
	dial := func(ctx context.Context) (Client, error) {
		client, err := NewClient(address)
		if err != nil {
			return nil, err
		}

		clientPtr.Store(&client)
		return client, nil
	}

	pool := NewPool(4, dial)
	defer pool.Close()

	for range 100 {
		func() {
			client, err := pool.Get(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer client.Close()

			poolClient, ok := client.(poolClient)
			if !ok {
				t.Fatalf("got %T is wrong", client)
			}

			loadClient := clientPtr.Load()
			if loadClient == nil {
				t.Fatalf("load client is nil")
			}

			got := fmt.Sprintf("%p", poolClient.client)
			want := fmt.Sprintf("%p", *loadClient)
			if got != want {
				t.Fatalf("got %s != want %s", got, want)
			}

			status := pool.Status()
			wantStatus := Status{Limit: 4, Using: 1, Idle: 0, Waiting: 0}
			if status != wantStatus {
				t.Fatalf("got %+v != want %+v", status, wantStatus)
			}
		}()
	}

	status := pool.Status()
	wantStatus := Status{Limit: 4, Using: 0, Idle: 1, Waiting: 0}
	if status != wantStatus {
		t.Fatalf("got %+v != want %+v", status, wantStatus)
	}

	if err = pool.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = pool.Get(ctx)
	if err != errPoolClosed {
		t.Fatalf("got %+v != want %+v", err, errPoolClosed)
	}
}
