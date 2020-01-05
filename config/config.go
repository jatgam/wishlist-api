package config

import (
	"github.com/jatgam/wishlist-api/utils"
)

type Config struct {
	DB           *DBConfig
	Secret       string
	JWTRealmName string
	EMail        *SGEmailConfig
}

type SGEmailConfig struct {
	SendGridAPIKey string
	FromName       string
	FromAddress    string
	Debug          bool
}

type DBConfig struct {
	Hostname string
	Database string
	User     string
	Password string
}

func GetConfig() *Config {
	var newConf Config

	newConf = Config{
		DB: &DBConfig{
			Hostname: utils.GetEnv("DB_HOSTNAME", "localhost"),
			Database: utils.GetEnv("DB_NAME", "wishlist"),
			User:     utils.GetEnv("DB_USERNAME", "wishlist"),
			Password: utils.GetEnv("DB_PASSWORD", "changeme"),
		},
		Secret:       utils.GetEnv("AUTH_SECRET", "A super secret key for jwt auth"),
		JWTRealmName: utils.GetEnv("JWT_REALM_NAME", "jatgam-wishlist"),
		EMail: &SGEmailConfig{
			SendGridAPIKey: utils.GetEnv("SENDGRID_API_KEY", ""),
			FromName:       utils.GetEnv("EMAIL_FROM_NAME", "Wishlist Admin"),
			FromAddress:    utils.GetEnv("EMAIL_FROM_ADDRESS", "wishlist@example.com"),
			Debug:          utils.GetEnvAsBool("EMAIL_DEBUG", "false"),
		},
	}

	return &newConf
}
