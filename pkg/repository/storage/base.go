package storage

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Database struct {
	Postgresql *gorm.DB
	Redis      *redis.Client
}

var DB *Database = &Database{}

func Connection() *Database {
	return DB
}
