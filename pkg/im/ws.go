package im

import (
	"fmt"
	"go-im/pkg/mq"
	"strings"
	"time"

	"gopkg.in/olahol/melody.v1"
)

var (
	// WS 实例
	Melody *melody.Melody

	// 用户与 session 相互映射
	userToSession = make(map[string]*melody.Session)
	sessionToUser = make(map[*melody.Session]string)
)

func init() {
	Melody = melody.New()

	// ws 处理连接
	Melody.HandleConnect(func(s *melody.Session) {
		fmt.Println("online")
		sessionToUser[s] = "" // 待确认 userID 的 session
		_ = s.Write([]byte("ping"))
	})

	// ws 处理断开
	Melody.HandleDisconnect(func(s *melody.Session) {
		// 清除绑定关系
		fmt.Println("offline")
		account := sessionToUser[s]
		if account != "" {
			delete(userToSession, account)
		}
		delete(sessionToUser, s)
	})

	// ws 处理消息
	Melody.HandleMessage(func(s *melody.Session, msg []byte) {
		raw := strings.Split(string(msg), ":")
		account := raw[0]
		message := raw[1]
		if message == "ping" {
			audience := []*melody.Session{s}
			// 绑定自己和群组的路由接收消息
			messages := mq.Receive(account, []string{"default_group"})
			for index := range messages {
				_ = Melody.BroadcastMultiple([]byte(messages[index]), audience)
			}
		} else if message == "pong" {
			// 绑定关系
			sessionToUser[s] = account
			userToSession[account] = s
		}
	})

	// 每 10 分钟清理一次失效连接
	go func() {
		timeTicker := time.NewTicker(time.Second * 60 * 10)
		for {
			for index := range sessionToUser {
				if index.IsClosed() || sessionToUser[index] == "" {
					_ = index.Close()
					if sessionToUser[index] != "" {
						delete(userToSession, sessionToUser[index])
					}
					delete(sessionToUser, index)
				}
			}
			<-timeTicker.C
		}
	}()
}
