package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.User{},
		models.AccessToken{},
		models.Room{},
		models.Profile{},
		models.UserRoom{},
		models.Message{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
