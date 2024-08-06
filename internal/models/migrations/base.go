package migrations

import (
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"gorm.io/gorm"
)

func RunAllMigrations(db *storage.Database) {

	// verification migration
	MigrateModels(db.Postgresql, AuthMigrationModels(), AlterColumnModels())

}

func MigrateModels(db *gorm.DB, models []interface{}, AlterColums []AlterColumn) {
	_ = db.AutoMigrate(models...)

	for _, d := range AlterColums {
		err := d.UpdateColumnType(db)
		if err != nil {
			fmt.Println("error migrating ", d.TableName, "for column", d.Column, ": ", err)
		}

	}

}
