package handler

import (
	"go-im/pkg/mq"

	"github.com/gin-gonic/gin"
)

func PushGroup(c *gin.Context) {
	group := c.DefaultPostForm("group", "default_group")
	msg := c.DefaultPostForm("msg", "")

	if group == "" || msg == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "参数缺失",
			"data":    "",
		})
		return
	}

	mq.Send(group, msg)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "已妥投",
		"data": map[string]string{
			"group": group,
			"msg":   msg,
		},
	})
}
