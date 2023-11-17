package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

func stressServer(address string, wg *sync.WaitGroup, data []byte) {
	defer wg.Done()

	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Println("Error connecting:", err)
			continue // Continue to the next iteration to retry connecting
		}

		_, err = conn.Write(data)
		if err != nil {
			log.Println("Error writing to connection:", err)
		}

		conn.Close()
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter IP: ")
	ip, _ := reader.ReadString('\n')
	ip = ip[:len(ip)-1] // Remove newline character

	fmt.Print("Enter port: ")
	port, _ := reader.ReadString('\n')
	port = port[:len(port)-1] // Remove newline character

	fmt.Print("Enter amount of threads: ")
	threadsStr, _ := reader.ReadString('\n')
	threadsStr = threadsStr[:len(threadsStr)-1] // Remove newline character
	threads, err := strconv.Atoi(threadsStr)
	if err != nil {
		log.Fatal(err)
	}

	address := ip + ":" + port
	data := make([]byte, 1024*1024) // 1MB
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go stressServer(address, &wg, data) // Start a new goroutine
	}

	wg.Wait()
}
