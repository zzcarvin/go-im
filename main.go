package main

import (
	"fmt"
	"net/http"

	"go-im/pkg/mq"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "view/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		s.Write([]byte("hi"))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		s.Write([]byte("bye"))
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println(string(msg))
		audience := []*melody.Session{s}
		messages := mq.Receive("123456", []string{})
		if len(messages) != 0 {
			for index := range messages {
				m.BroadcastMultiple([]byte(messages[index]), audience)
			}
		} else if string(msg) == "ping" {
			m.BroadcastMultiple([]byte("pong"), audience)
		} else {
			m.BroadcastMultiple(msg, audience)
		}
	})

	r.Run(":5000")
}
