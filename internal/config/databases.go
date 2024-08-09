package config

type Database struct {
	DB_HOST       string
	DB_PORT       string
	DB_CONNECTION string
	USERNAME      string
	PASSWORD      string
	TIMEZONE      string
	SSLMODE       string
	DB_NAME       string
	Migrate       bool
}

type Redis struct {
	REDIS_PORT string
	REDIS_HOST string
	REDIS_DB   string
}
