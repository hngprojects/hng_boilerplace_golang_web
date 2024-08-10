package models

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	dbRedis "github.com/hngprojects/telex_be/pkg/repository/storage/redis"
)

type NotificationRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Sent bool   `json:"sent"`
}

type SendOTP struct {
	Email    string `json:"email"  validate:"required"`
	OtpToken int    `json:"otp_token"  validate:"required"`
}

type SendWelcomeMail struct {
	Email string `json:"email"  validate:"required"`
}

type SendEmailVerificationMail struct {
	Email string `json:"email"`
	Code  uint   `json:"code" validate:"required"`
	Token string `json:"token"`
}

type SendResetPassword struct {
	Email string `json:"email"`
	Token int    `json:"token"  validate:"required"`
}

type SendMagicLink struct {
	Email     string `json:"email"  validate:"required"`
	MagicLink string `json:"magic_link"  validate:"required"`
}

type SendSqueeze struct {
	Email     string `json:"email"  validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
type SendContactUsMail struct {
	Name    string `json:"name"  validate:"required"`
	Email   string `json:"email" `
	Subject string `son:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func (n *NotificationRecord) PushToQueue(rdb *redis.Client) error {
	err := dbRedis.PushToQueue(rdb, &n)

	if err != nil {
		return err
	}

	return nil
}

func (n *NotificationRecord) PopFromQueue(rdb *redis.Client) (NotificationRecord, error) {
	var rec NotificationRecord
	res, err := dbRedis.PopFromQueue(rdb)

	if err != nil {
		return rec, err
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		return rec, fmt.Errorf("error marshaling map: %v", err)
	}

	err = json.Unmarshal(resJSON, &rec)
	if err != nil {
		return rec, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return rec, nil
}
