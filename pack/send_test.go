// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/FishGoddess/vex"
)

func runTestServer(t *testing.T, address string, str string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Error(err)
	}

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	buf := make([]byte, headerSize+len(str))

	n, err := io.ReadFull(conn, buf[:])
	if err != nil {
		t.Error(err)
	}

	if n < headerSize {
		t.Errorf("n %d < headerSize %d", n, headerSize)
	}

	header := Endian.Uint64(buf[:headerSize])
	magic := (header >> (typeBits + dataSizeBits)) & maxMagic
	if magic != magicNumber {
		t.Errorf("magic %d != magicNumber %d", magic, magicNumber)
	}

	packetType := byte((header >> dataSizeBits) & maxType)
	dataSize := int32(header & maxDataSize)
	if packetType != packetTypeTest {
		t.Errorf("packetType %d != packetTypeTest %d", packetType, packetTypeTest)
	}

	data := buf[headerSize : headerSize+dataSize]
	if string(data) != str {
		t.Errorf("data %s is wrong!", data)
	}

	data = []byte(str)
	header = magicNumber<<(typeBits+dataSizeBits) | uint64(packetTypeStandard)<<dataSizeBits | uint64(len(data))

	var headerBytes [headerSize]byte
	Endian.PutUint64(headerBytes[:], header)

	n, err = conn.Write(headerBytes[:])
	if err != nil {
		t.Error(err)
	}

	if n != headerSize {
		t.Errorf("n %d != headerSize %d", n, headerSize)
	}

	n, err = conn.Write(data)
	if err != nil {
		t.Error(err)
	}

	if n != int(dataSize) {
		t.Errorf("n %d != dataSize %d", n, dataSize)
	}
}

// go test -v -cover -run=^TestSend$
func TestSend(t *testing.T) {
	address := "127.0.0.1:9988"

	str := "key value"
	go runTestServer(t, address, str)
	time.Sleep(time.Second)

	client, err := vex.NewClient(address)
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	packet, err := Send(client, packetTypeTest, []byte(str))
	if err != nil {
		t.Error(err)
	}

	if string(packet) != str {
		t.Errorf("packet %s is wrong!", packet)
	}
}
