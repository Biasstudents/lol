package main

import (
  "bufio"
  "fmt"
  "log"
  "math/rand"
  "net"
  "os"
  "strconv"
  "sync"
  "time"
)

var totalDataSent int64 = 0
var debug = false // Set this to true if you want to see error messages

func generateMalformedData() []byte {
  // Generate random data with sizes between 1MB and 10MB
  size := rand.Intn(10*1024*1024-1*1024*1024+1) + 1*1024*1024
  randomData := make([]byte, size)
  rand.Read(randomData)

  // Introduce intentional errors in the data
  for i := 0; i < size*9/10; i++ { // Introduce errors in 40% of the data
    randomIndex := rand.Intn(len(randomData))
    randomData[randomIndex] = byte(rand.Intn(256)) // Modify a byte to introduce error
  }

  return randomData
}



func stressServer(address string, wg *sync.WaitGroup, maxConnections int) {
  defer wg.Done()

  conn, err := net.Dial("tcp", address)
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  for {
    malformedData := generateMalformedData()

    _, err := conn.Write(malformedData)
    if err != nil && debug {
      log.Println("Error writing to connection:", err)
    }

    totalDataSent += int64(len(malformedData))

  }
}

func printBandwidth() {
  ticker := time.NewTicker(1 * time.Second)
  defer ticker.Stop()

  for range ticker.C {
    switch {
    case totalDataSent > 1<<30:
      fmt.Printf("\rTotal data sent: %.2f GB", float64(totalDataSent)/(1<<30))
    case totalDataSent > 1<<20:
      fmt.Printf("\rTotal data sent: %.2f MB", float64(totalDataSent)/(1<<20))
    case totalDataSent > 1<<10:
      fmt.Printf("\rTotal data sent: %.2f KB", float64(totalDataSent)/(1<<10))
    default:
      fmt.Printf("\rTotal data sent: %d bytes", totalDataSent)
    }
    totalDataSent = 0
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

  address := ip + ":" + port

  var wg sync.WaitGroup
  for i := 0; i < threads; i++ {
    wg.Add(1)
    go stressServer(address, &wg, 10) // 10 connections per thread
  }

  go printBandwidth()

  wg.Wait()
}
