package migrations

import "github.com/hngprojects/telex_be/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.User{},
		models.AccessToken{},
		models.Room{},
		models.Profile{},
		models.UserRoom{},
		models.Message{},
		models.MagicLink{},
		models.PasswordReset{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
