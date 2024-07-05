package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "os/signal"
    "runtime"
    "strconv"
    "sync"
    "sync/atomic"
    "syscall"
    "time"

    "github.com/valyala/fasthttp"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    // Ask for the target URL
    fmt.Print("Enter the target URL: ")
    targetURL, _ := reader.ReadString('\n')
    targetURL = targetURL[:len(targetURL)-1] // Remove the newline character

    // Ask for the number of concurrent threads
    fmt.Print("Enter the number of concurrent threads: ")
    concurrencyStr, _ := reader.ReadString('\n')
    concurrencyStr = concurrencyStr[:len(concurrencyStr)-1] // Remove the newline character
    concurrency, err := strconv.Atoi(concurrencyStr)
    if err != nil {
        fmt.Println("Invalid number of concurrent threads. Please enter a valid integer.")
        return
    }

    // Set GOMAXPROCS to use all available CPUs
    runtime.GOMAXPROCS(runtime.NumCPU())

    // Custom HTTP client with keep-alive and connection pooling
    client := &fasthttp.Client{
        MaxConnsPerHost:      concurrency * 2, // Allow more connections per host
        ReadTimeout:         2 * time.Second,  // Reduce timeouts
        WriteTimeout:        2 * time.Second,
        MaxIdleConnDuration: 10 * time.Second, // Keep connections alive longer
    }

    var wg sync.WaitGroup
    var requestCount int64

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    reqPool := sync.Pool{
        New: func() interface{} {
            return fasthttp.AcquireRequest()
        },
    }

    respPool := sync.Pool{
        New: func() interface{} {
            return fasthttp.AcquireResponse()
        },
    }

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                default:
                    req := reqPool.Get().(*fasthttp.Request)
                    resp := respPool.Get().(*fasthttp.Response)
                    req.SetRequestURI(targetURL)

                    if err := client.Do(req, resp); err != nil {
                        fmt.Println("Error:", err)
                    } else {
                        atomic.AddInt64(&requestCount, 1)
                    }

                    req.Reset()
                    resp.Reset()
                    reqPool.Put(req)
                    respPool.Put(resp)
                }
            }
        }()
    }

    go func() {
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        prevCount := int64(0)
        for range ticker.C {
            currentCount := atomic.LoadInt64(&requestCount)
            rps := (currentCount - prevCount) * 2 // as ticker interval is 0.5s
            prevCount = currentCount
            fmt.Printf("\rRPS: %d", rps)
        }
    }()

    // Graceful shutdown on Ctrl+C
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    cancel()
    wg.Wait()
}
