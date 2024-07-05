package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	bps     int64 = 0
	debug   bool  = false
	dataSize      = 1024 * 1024 // Default data size to 1MB (1,048,576 bytes)
	ip      string
	port    string
	threads int
)

func connectToServer(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func stressServer(ctx context.Context, address string, data []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := connectToServer(address)
			if err != nil {
				if debug {
					log.Printf("Error connecting: %v\n", err)
				}
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			n, err := conn.Write(data)
			if err != nil && debug {
				log.Printf("Error writing to connection: %v\n", err)
			}

			atomic.AddInt64(&bps, int64(n)) // Accumulate bytes sent
			conn.Close()
		}
	}
}

func printBandwidth(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			dataSent := atomic.SwapInt64(&bps, 0)
			bitsSent := dataSent * 8 // Convert bytes to bits
			bandwidth := float64(bitsSent) / elapsed
			startTime = time.Now()

			fmt.Printf("\rTotal data sent: ")
			switch {
			case bandwidth > 1<<30:
				fmt.Printf("%.2f Gbps", bandwidth/(1<<30))
			case bandwidth > 1<<20:
				fmt.Printf("%.2f Mbps", bandwidth/(1<<20))
			case bandwidth > 1<<10:
				fmt.Printf("%.2f Kbps", bandwidth/(1<<10))
			default:
				fmt.Printf("%d bps", int64(bandwidth))
			}
			fmt.Print("          ") // Ensure the line is cleared properly
		}
	}
}

func main() {
	flag.StringVar(&ip, "ip", "", "Server IP address")
	flag.StringVar(&port, "port", "", "Server port")
	flag.IntVar(&threads, "threads", 1, "Number of threads")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	if ip == "" || port == "" {
		log.Fatal("IP and port must be specified")
	}

	address := ip + ":" + port
	data := make([]byte, dataSize)
	ctx, cancel := context.WithCancel(context.Background())

	// Initial connection check
	var conn net.Conn
	var err error
	for {
		conn, err = connectToServer(address)
		if err == nil {
			conn.Close()
			fmt.Println("Stress testing started")
			break
		} else {
			if debug {
				log.Printf("Initial connection failed: %v\n", err)
			}
			// No delay here
		}
	}

	// Start stress testing threads
	for i := 0; i < threads; i++ {
		go stressServer(ctx, address, data)
	}

	go printBandwidth(ctx)

	// Handle termination signal to gracefully stop goroutines
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
	fmt.Println("\nExiting...")
}
