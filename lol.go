package main

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"log"
	"sync"
)

type exampleClient struct {
	data []byte
}

func (s *exampleClient) OnConnect(c *connection.Connection) {
	log.Println("OnConnect")
	c.Send(s.data)
}

func (s *exampleClient) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	log.Println("OnMessage")
	return
}

func (s *exampleClient) OnClose(c *connection.Connection) {
	log.Println("OnClose")
}

func main() {
	data := make([]byte, 1024*1024) // 1 MB of data
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := gev.NewClient("193.228.196.49:80", &exampleClient{data: data}, nil)
			if err != nil {
				log.Fatalln("NewClient failed:", err)
			}

			for {
				client.Send(data)
			}
		}()
	}

	wg.Wait()
}
