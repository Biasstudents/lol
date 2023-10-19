package main

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/eventloop"
	"log"
	"sync"
)

type exampleClient struct {
	data []byte
}

func (s *exampleClient) OnInitComplete(c *gev.Connection) {
	log.Println("OnInitComplete")
	c.Send(s.data)
}

func (s *exampleClient) OnMessage(c *gev.Connection, data []byte) {
	log.Println("OnMessage")
}

func (s *exampleClient) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func main() {
	data := make([]byte, 1024*1024) // 1 MB of data
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loop := eventloop.New()
			client, err := gev.NewClient(loop, "193.228.196.49:80", &exampleClient{data: data}, nil)
			if err != nil {
				log.Fatalln("NewClient failed:", err)
			}
			loop.Run()
		}()
	}

	wg.Wait()
}
