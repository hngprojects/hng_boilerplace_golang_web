package utility

import (
	"net/mail"
	"os"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/nyaruka/phonenumbers"
)

func EmailValid(email string) (string, bool) {
	// made some change to parse the formated email
	e, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return e.Address, err == nil
}

func PhoneValid(phone string) (string, bool) {
	parsed, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return phone, false
	}

	if !phonenumbers.IsValidNumber(parsed) {
		return phone, false
	}

	formattedNum := phonenumbers.Format(parsed, phonenumbers.NATIONAL)
	return formattedNum, true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CleanStringInput(input string) string {
	policy := bluemonday.UGCPolicy()
	cleanedInput := policy.Sanitize(input)
	re := regexp.MustCompile(`[^\w\s]`)
	cleanedInput = re.ReplaceAllString(cleanedInput, "")

	return cleanedInput
}
