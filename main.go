package main

import (
	"net/http"

	"go-im/internal/handler"
	"go-im/pkg/im"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// H5 客户端
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "view/index.html")
	})

	// WebSocket 服务
	r.GET("/ws", func(c *gin.Context) {
		im.Melody.HandleRequest(c.Writer, c.Request)
	})

	// 推送个人
	r.POST("/push/user", handler.PushUser)

	// 推送群组
	r.POST("/push/group", handler.PushGroup)

	r.Run(":5000")
}
