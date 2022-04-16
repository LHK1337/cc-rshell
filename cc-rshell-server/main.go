package main

import (
	"cc-rshell-server/sockets"
	"github.com/gin-gonic/gin"
	"log"
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
	log.Panicln(engine.Run())
}
