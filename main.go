package main

import (
	"net/http"

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
		audience := []*melody.Session{s}
		m.BroadcastMultiple(msg, audience)
	})

	r.Run(":5000")
}
