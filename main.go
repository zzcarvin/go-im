package main

import (
	"fmt"
	"net/http"
	"strings"

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
		// 从在线连接中注册 undo
		fmt.Println("bind")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		// 从在线连接中清除 undo
		fmt.Println("unbind")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		raw := strings.Split(string(msg), ":")
		account := raw[0]
		message := raw[1]
		if message == "ping" {
			audience := []*melody.Session{s}
			messages := mq.Receive(account, []string{})
			for index := range messages {
				m.BroadcastMultiple([]byte(messages[index]), audience)
			}
		} else if message == "pong" {
			// 更新连接活跃度 undo
			fmt.Println("refresh")
		}
	})

	r.Run(":5000")
}
