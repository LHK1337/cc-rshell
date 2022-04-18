package main

import (
	"cc-rshell-server/sockets"
	"cc-rshell-server/ssh"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func main() {
	engine := gin.Default()
	clients := sockets.NewClientSocketHandler()

	clientsRouter := engine.Group("/clients")
	clientsRouter.GET("/socket", func(c *gin.Context) {
		err := clients.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{})
		if err != nil {
			log.Panicln(err)
		}
	})

	defer clients.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Panicln(engine.Run())
	}()

	go func() {
		defer wg.Done()
		log.Panicln(ssh.ListenAndServer("127.0.0.1:2222", clients.ClientRegistry))
	}()

	wg.Wait()
}
