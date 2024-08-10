package notifications

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/send"
)

func (n NotificationObject) SendContactUsMail() error {
	var (
		notificationData     = models.SendContactUsMail{}
		subject              = "Contact us message received"
		templateFileName     = "default.html"
		baseTemplateFileName = ""
	)

	err := json.Unmarshal([]byte(n.Notification.Data), &notificationData)
	if err != nil {
		return fmt.Errorf("error decoding saved notification data, %v", err)
	}

	data, err := ConvertToMapAndAddExtraData(notificationData, map[string]interface{}{"firstname": thisOrThatStr(notificationData.Name, notificationData.Email), "business_name": thisOrThatStr("", "")})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v, %v", err, strings.Join([]string{err.Error()}, ", "))
	}

	return send.SendEmail(n.ExtReq, notificationData.Email, subject, templateFileName, baseTemplateFileName, data)
}
