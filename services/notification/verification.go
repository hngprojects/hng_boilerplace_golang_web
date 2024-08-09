package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/send"
)

func (n NotificationObject) SendEmailVerificationMail() error {
	var (
		notificationData     = models.SendEmailVerificationMail{}
		subject              = "Subject: Please verify your email address"
		templateFileName     = "email_verification.html"
		baseTemplateFileName = ""
		configData           = config.GetConfig()
		user                 models.User
	)

	err := json.Unmarshal([]byte(n.Notification.Data), &notificationData)
	if err != nil {
		return fmt.Errorf("error decoding saved notification data, %v", err)
	}

	user, err = user.GetUserByEmail(n.Db, notificationData.Email)
	if err != nil {
		return fmt.Errorf("error getting user with account id %v, %v", notificationData.Email, err)
	}

	verificationUrl := fmt.Sprintf("%v/email-verify/", configData.App.Url)

	data, err := ConvertToMapAndAddExtraData(notificationData, map[string]interface{}{"firstname": thisOrThatStr(user.Profile.FirstName, user.Email), "verification_url": verificationUrl})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v", err)
	}

	return send.SendEmail(n.ExtReq, user.Email, subject, templateFileName, baseTemplateFileName, data)
}
