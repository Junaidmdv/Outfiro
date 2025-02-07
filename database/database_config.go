package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Connect struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     string
	Ssl_mode string
	TimeZone string
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("this is env error load")
	}

}

func DbConnection() {

	DbConfig := Connect{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Ssl_mode: os.Getenv("SSL_MODE"),
		TimeZone: os.Getenv("TIME_ZONE"),
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		DbConfig.Host, DbConfig.User, DbConfig.Password, DbConfig.Dbname, DbConfig.Port, DbConfig.Ssl_mode, DbConfig.TimeZone,
	)
	fmt.Printf("%s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("database connection error")
	}

	DB = db
}
