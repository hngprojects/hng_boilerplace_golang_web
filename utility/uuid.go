package utility

import "github.com/gofrs/uuid"

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func IsValidUUID(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}
