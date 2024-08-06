package storage

import "gorm.io/gorm"

type Database struct {
	Postgresql *gorm.DB
}

var DB *Database = &Database{}

func Connection() *Database {
	return DB
}
