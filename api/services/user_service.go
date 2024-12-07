package services

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type User struct {
	Name       string           `json:"name,omitempty"`
	TotalPoint int              `json:"value,omitempty"`
	nofication chan interface{} `json:"-"`

	Conn *websocket.Conn `json:"-"`
	quiz *QuizSession    `json:"-"`
}

func NewUser(quiz *QuizSession, conn *websocket.Conn, name string) *User {
	user := &User{
		Name:       name,
		nofication: make(chan interface{}),
		quiz:       quiz,
		Conn:       conn,
	}

	go func() {
		for message := range user.nofication {
			if err := user.Conn.WriteJSON(message); err != nil {
				log.Err(err).Msg("can not send message")
			}
			time.Sleep(time.Millisecond)
		}
	}()

	quiz.Join(user)
	return user
}

type MessageType string

const (
	Echo   MessageType = "echo"
	Sync   MessageType = "sync"
	Update MessageType = "update"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

func (u *User) ResetPoint() {
	u.TotalPoint = 0
}

// func (u *User) Listen(message Message) {
// 	if u.Conn == nil {
// 		return
// 	}

// 	if message.Type == Echo {
// 		if msg := marshalData(message.Data); msg != nil {
// 			if err := u.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
// 				log.Warn().AnErr("user listen", err)
// 			}
// 		}
// 		return
// 	}

// 	msg, err := json.Marshal(message)
// 	if err != nil {
// 		log.Warn().AnErr("user listen: marshal message: ", err)
// 		return
// 	}

// 	if err := u.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
// 		log.Warn().AnErr("user listen", err)
// 	}
// }
