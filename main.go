package main

import (
	"fmt"
	"net/http"
	"strings"

	"go-im/internal/handler"
	"go-im/pkg/mq"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	// H5 客户端
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "view/index.html")
	})

	// WebSocket 服务
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// 推送个人
	r.POST("/push/user", handler.PushUser)

	// 推送群组
	r.POST("/push/group", handler.PushGroup)

	// ws 处理连接
	m.HandleConnect(func(s *melody.Session) {
		// 从在线连接中注册 undo
		fmt.Println("bind")
	})

	// ws 处理断开
	m.HandleDisconnect(func(s *melody.Session) {
		// 从在线连接中清除 undo
		fmt.Println("unbind")
	})

	// ws 处理消息
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		raw := strings.Split(string(msg), ":")
		account := raw[0]
		message := raw[1]
		if message == "ping" {
			audience := []*melody.Session{s}
			// 绑定自己和群组的路由接收消息
			messages := mq.Receive(account, []string{"default_group"})
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
