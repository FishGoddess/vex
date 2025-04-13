// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/FishGoddess/vex"
)

func runClient(msg string) {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	if _, err := client.Write([]byte(msg)); err != nil {
		panic(err)
	}

	var buf [1024]byte
	n, err := client.Read(buf[:])
	if err != nil {
		panic(err)
	}

	fmt.Println("Received:", string(buf[:n]))
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		msg := strconv.Itoa(i)

		wg.Add(1)
		go func() {
			defer wg.Done()

			runClient(msg)
		}()
	}

	wg.Wait()
	fmt.Println("run clients done.")
}
