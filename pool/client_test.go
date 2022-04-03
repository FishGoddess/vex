// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import "testing"

// go test -v -cover -run=^TestWrapClient$
func TestWrapClient(t *testing.T) {
	pool := &Pool{}
	client := wrapClient(pool, nil)

	poolClient, ok := client.(*poolClient)
	if !ok {
		t.Errorf("client.(*poolClient) %T not ok", client)
	}

	if poolClient.pool != pool {
		t.Errorf("poolClient.pool %p != pool %p", poolClient.pool, pool)
	}

	if poolClient.client != nil {
		t.Errorf("poolClient.client %p != nil", poolClient.client)
	}
}
