package pool

import (
	"errors"
	"sync"

	"github.com/FishGoddess/vex"
)

var (
	errClientClosed = errors.New("vex: client is closed")
)

// poolClient wraps client to a pool client.
type poolClient struct {
	pool   *Pool
	client vex.Client
	closed bool
	lock   sync.RWMutex
}

// wrapClient wraps client to a pool client.
func wrapClient(pool *Pool, client vex.Client) vex.Client {
	return &poolClient{
		pool:   pool,
		client: client,
		closed: false,
	}
}

// Send sends a packet with requestBody to server and returns responseBody responded from server.
func (pc *poolClient) Send(packetType vex.PacketType, requestBody []byte) (responseBody []byte, err error) {
	pc.lock.RLock()
	defer pc.lock.RUnlock()

	if pc.closed {
		return nil, errClientClosed
	}

	return pc.client.Send(packetType, requestBody)
}

// Close closes current client.
func (pc *poolClient) Close() error {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	if pc.closed {
		return nil
	}

	pc.pool.put(pc)
	return nil
}
