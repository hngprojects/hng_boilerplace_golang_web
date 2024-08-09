package actions

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/actions/names"
	notifications "github.com/hngprojects/telex_be/services/notification"
)

func Send(extReq request.ExternalRequest, db *gorm.DB, rdb *redis.Client, notification *models.NotificationRecord) error {
	var (
		err  error
		req  = notifications.NewNotificationObject(extReq, rdb, db, notification)
		name = GetName(notification.Name)
	)

	callEmailFunc := map[names.NotificationName]func() error{
		names.SendOTP: func() error {
			return req.SendOTP()
		},
		names.SendWelcomeMail: func() error {
			return req.SendWelcomeMail()
		},
		names.SendResetPasswordMail: func() error {
			return req.SendResetPasswordMail()
		},
		names.SendEmailVerificationMail: func() error {
			return req.SendEmailVerificationMail()
		},
		names.SendMagicLink: func() error {
			return req.SendMagicLink()
		},
		names.SendSqueeze: func() error {
			return req.SendSqueeze()
		},
		names.SendContactUsMail: func() error {
			return req.SendContactUsMail()
		},
	}

	err = callEmailFunc[name]()

	if err != nil {
		return err
	}

	return nil
}

func GetName(name string) names.NotificationName {
	return names.NotificationName(name)
}
