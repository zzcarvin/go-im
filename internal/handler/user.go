package handler

import (
	"go-im/pkg/mq"

	"github.com/gin-gonic/gin"
)

func PushUser(c *gin.Context) {
	user := c.DefaultPostForm("user", "")
	msg := c.DefaultPostForm("msg", "")

	if user == "" || msg == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "参数缺失",
			"data":    "",
		})
		return
	}

	mq.Send(user, msg)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "已妥投",
		"data": map[string]string{
			"user": user,
			"msg":  msg,
		},
	})
}
