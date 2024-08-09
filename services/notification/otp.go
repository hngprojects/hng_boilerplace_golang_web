package notifications

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/services/send"
)

func SendDirectOTP(extReq request.ExternalRequest, req models.SendOTP, db *gorm.DB) error {
	var (
		templateFileName     = "send_otp.html"
		baseTemplateFileName = ""
		errs                 []string
		user                 models.User
	)

	subject := fmt.Sprintf("Subject: Secure Login: Your OTP Code Is: %v", req.OtpToken)

	user, err := user.GetUserByEmail(db, req.Email)
	if err != nil {
		return fmt.Errorf("error getting user with account id %v, %v", req.Email, err)
	}

	data, err := ConvertToMapAndAddExtraData(req, map[string]interface{}{"firstname": thisOrThatStr(user.Profile.FirstName, user.Email), "business_name": thisOrThatStr("", "Telex")})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v, %v", err, strings.Join(errs, ", "))
	}

	err = send.SendEmail(extReq, user.Email, subject, templateFileName, baseTemplateFileName, data)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}
	return nil
}

func (n NotificationObject) SendOTP() error {
	var (
		notificationData     = models.SendOTP{}
		templateFileName     = "send_otp.html"
		baseTemplateFileName = ""
		errs                 []string
		user                 models.User
	)

	err := json.Unmarshal([]byte(n.Notification.Data), &notificationData)
	if err != nil {
		return fmt.Errorf("error decoding saved notification data, %v", err)
	}

	subject := fmt.Sprintf("Subject: Secure Login: Your OTP Code Is: %v", notificationData.OtpToken)

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
