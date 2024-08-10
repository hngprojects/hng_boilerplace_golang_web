package names

import (
	"fmt"
	"reflect"

	"github.com/hngprojects/telex_be/utility"
)

type NotificationName string

const (
	SendWelcomeMail           NotificationName = "send_welcome_mail"
	SendOTP                   NotificationName = "send_otp"
	SendResetPasswordMail     NotificationName = "send_reset_password_mail"
	SendEmailVerificationMail NotificationName = "send_email_verification_mail"
	SendMagicLink             NotificationName = "send_magic_link"
	SendSqueeze               NotificationName = "send_squeeze"
	SendContactUsMail         NotificationName = "send_contact_us"
)

func Check() {
	constantName := "SendWelcomeMail"
	constantValue := reflect.ValueOf(constantName).Interface().(string)
	fmt.Println("check", constantValue)
}

func GetNames(pkgImportPath string) ([]string, error) {
	names := []string{}
	constants, err := utility.GetConstants(pkgImportPath)
	if err != nil {
		return names, err
	}

	for _, v := range constants {
		names = append(names, v)
	}

	return names, nil
}
