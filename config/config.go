package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	ENVIRONMENT     = ""
	PORT            = ""
	MIGRATE         = false
	DB_URI          = ""
	REDIS_DB_NUMBER = ""
)

func LoadConfig() {
	ENVIRONMENT = viper.GetString("ENVIRONMENT")
	PORT = viper.GetString("PORT")
	MIGRATE = viper.GetBool("MIGRATE")
	REDIS_DB_NUMBER = viper.GetString("REDIS_DB_NUMBER")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbName := viper.GetString("DB_NAME")
	dbPassword := viper.GetString("DB_PASSWORD")

	DB_URI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", dbHost, dbUser, dbPassword, dbName, dbPort)
}
