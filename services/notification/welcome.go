package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/send"
)

func (n NotificationObject) SendWelcomeMail() error {
	var (
		notificationData     = models.SendWelcomeMail{}
		subject              = "Subject: Welcome on board!🎉"
		templateFileName     = "welcome-email.html"
		baseTemplateFileName = ""
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

	data, err := ConvertToMapAndAddExtraData(notificationData, map[string]interface{}{"firstname": thisOrThatStr(user.Profile.FirstName, user.Email)})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v", err)
	}

	return send.SendEmail(n.ExtReq, user.Email, subject, templateFileName, baseTemplateFileName, data)
}
