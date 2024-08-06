package utility

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

func GetRandomNumbersInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomString(length int) string {
	u, _ := uuid.NewV4()
	uuidStr := u.String()
	// Regular expression pattern to match all non-alphanumeric characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return ""
	}
	// Replacing all non-alphanumeric characters with an empty string
	processedString := reg.ReplaceAllString(uuidStr+uuidStr[:length%36], "")
	if len(processedString) >= length {
		return processedString[:length]
	}
	// Padding the processed string with random alphanumeric characters to make it the desired length
	alphanumeric := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	padding := make([]byte, length-len(processedString))
	rand.Read(padding)
	for i, b := range padding {
		padding[i] = alphanumeric[b%byte(len(alphanumeric))]
	}
	return processedString + string(padding)
}
