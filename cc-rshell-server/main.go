package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"time"
)

func main() {
	engine := gin.Default()
	clients := melody.New()
	// TODO: better origin check
	clients.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	engine.GET("/", func(c *gin.Context) {
		err := clients.HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Panicln(err)
		}
	})

	go func() {
		for {
			bytes, _ := json.Marshal(gin.H{
				"type":   "event",
				"event":  "char",
				"params": []string{"a"},
			})
			time.Sleep(time.Millisecond * 100)

			clients.Broadcast(bytes)
			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "char",
				"params": []string{"b"},
			})
			time.Sleep(time.Millisecond * 100)

			clients.Broadcast(bytes)
			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "char",
				"params": []string{"c"},
			})
			time.Sleep(time.Millisecond * 100)

			clients.Broadcast(bytes)
			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "char",
				"params": []string{"d"},
			})
			time.Sleep(time.Millisecond * 100)

			clients.Broadcast(bytes)
			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "char",
				"params": []string{"e"},
			})
			time.Sleep(time.Millisecond * 100)

			clients.Broadcast(bytes)

			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "key",
				"params": []interface{}{257, false},
			})
			clients.Broadcast(bytes)

			bytes, _ = json.Marshal(gin.H{
				"type":   "event",
				"event":  "key_up",
				"params": []interface{}{257},
			})
			clients.Broadcast(bytes)

			time.Sleep(time.Second * 1)
			bytes, _ = json.Marshal(gin.H{
				"type":    "serverNotification",
				"message": "done!",
			})
			clients.Broadcast(bytes)
			time.Sleep(time.Second * 5)
		}
	}()

	clients.HandleMessage(func(s *melody.Session, bytes []byte) {

	})

	defer func(clients *melody.Melody) {
		_ = clients.Close()
	}(clients)
	log.Panicln(engine.Run())
}
