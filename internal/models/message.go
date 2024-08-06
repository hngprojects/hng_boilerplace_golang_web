package models

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type Message struct {
	ID        int       `gorm:"column:id; type:serial; primaryKey" json:"id"`
	Content   string    `gorm:"column:content; type:text; not null" json:"content"`
	RoomID    string    `gorm:"type:uuid;not null" json:"room_id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	Username  string    `gorm:"column:username; type:varchar(255)" json:"username"`
	CreatedAt time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
}

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required"`
	UserId  string `json:"user_id"`
	RoomId  string `json:"room_id"`
}

func (m *Message) CreateMessage(db *gorm.DB) error {

	var userRoom UserRoom

	exist := postgresql.CheckExists(db, &userRoom, "room_id = ? AND user_id = ?", m.RoomID, m.UserID)
	if !exist {
		return errors.New("user not in room")
	}

	m.Username = userRoom.Username

	err := postgresql.CreateOneRecord(db, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) GetMessagesByRoomID(db *gorm.DB, userId, roomID string) ([]Message, error) {
	var messages []Message
	var userRoom UserRoom

	exist := postgresql.CheckExists(db, &userRoom, "room_id = ? AND user_id = ?", roomID, userId)
	if !exist {
		return messages, errors.New("user not in room")
	}

	err := postgresql.SelectAllFromDb(db, "", &messages, "room_id = ?", roomID)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

func (m *Message) GetMessageByID(db *gorm.DB, messageID string) (Message, error) {
	var message Message

	err, nerr := postgresql.SelectOneFromDb(db, &message, "message_id = ?", messageID)
	if err != nil {
		return message, nerr
	}
	return message, nil
}
