package models

import (
	"outfiro/database"
	"outfiro/utils"
)

func Migrate() {
	database.DB.AutoMigrate(&Users{})
	database.DB.AutoMigrate(&Products{})
	database.DB.AutoMigrate(&Categories{})
	database.DB.AutoMigrate(&Admin{})
	database.DB.AutoMigrate(&utils.OTP{})
}