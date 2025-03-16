package models

import (
	"outfiro/database"
)

func Migrate() {
	database.DB.AutoMigrate(&Users{})
	database.DB.AutoMigrate(&Products{})
	database.DB.AutoMigrate(&Review{})
	database.DB.AutoMigrate(&Categories{})
	database.DB.AutoMigrate(&Admin{})
	database.DB.AutoMigrate(&OtpRecord{})
	database.DB.AutoMigrate(&SignuPlayload{})
	database.DB.AutoMigrate(&Address{})
	database.DB.AutoMigrate(&CartItems{})
	database.DB.AutoMigrate(&Order{})
	database.DB.AutoMigrate(&OrderItem{})
	database.DB.AutoMigrate(&ShippingAddress{})
	database.DB.AutoMigrate(&Wishlist_item{})
	database.DB.AutoMigrate(&Payment{})
	database.DB.AutoMigrate(&Coupon{})
	database.DB.AutoMigrate(&CouponLimit{})
	database.DB.AutoMigrate(&WalleteHistory{})
	database.DB.AutoMigrate(&Transaction{})
}
