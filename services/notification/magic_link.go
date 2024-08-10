package notifications

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/send"
)

func (n NotificationObject) SendMagicLink() error {
	var (
		notificationData     = models.SendMagicLink{}
		templateFileName     = "send_magic_link.html"
		baseTemplateFileName = ""
		errs                 []string
		user                 models.User
	)

	err := json.Unmarshal([]byte(n.Notification.Data), &notificationData)
	if err != nil {
		return fmt.Errorf("error decoding saved notification data, %v", err)
	}

	subject := "Subject: Secure Login: Your MagicLink..."

	user, err = user.GetUserByEmail(n.Db, notificationData.Email)
	if err != nil {
		return fmt.Errorf("error getting user with account id %v, %v", notificationData.Email, err)
	}

	data, err := ConvertToMapAndAddExtraData(notificationData, map[string]interface{}{"firstname": thisOrThatStr(user.Profile.FirstName, user.Email), "business_name": thisOrThatStr("", "")})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v", err)
	}

	err = send.SendEmail(n.ExtReq, user.Email, subject, templateFileName, baseTemplateFileName, data)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}
	return nil
}
