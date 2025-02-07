package models

import (
	"gorm.io/gorm"
	"time"
)

type Categories struct {
	gorm.Model
	CategoryName string     `json:"name" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	Products     []Products `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type Products struct {
	gorm.Model
	ProductName string  `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Description string  `json:"description" validate:"required,min=3,max=100,required"`
	ImageUrl    string   `json:"image_url"  validate:"required,url"`
	CategoryId  uint    `json:"categoryid" binding:"required"`
	Price       float64 `json:"price" binding:"required" validate:"numeric,gt=0,required "`
	Status      string  `json:"status" validate:"required oneof='In stock''Out of stokck'"`
	Size        string  `json:"size" binding:"required"`
	Quntity     uint    `json:"quantity" validate:"required,numeric,gt=0"`
}
type Users struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname" validate:"required,alpha_space"`
	LastName  string `gorm:"not null" json:"lastname"  validate:"required,alpha_space"`
	Email     string `gorm:"not null;unique_index" json:"email" validate:"email"`
	Address   string `json:"address" validate:"omitempty"`
	Password  string `json:"password"`
	Role      string `gorm:"default:user" json:"role"`
	IsBlocked bool   `gorm:"default:false" json:"isblocked"`
	LoginType string `gorm:"default:normal"`
}
type Admin struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname" validate:"required,alpha_space"`
	LastName  string `gorm:"not null" json:"lastname"  validate:"omitempty"`
	Email     string `gorm:"not null;unique_index" json:"email"  validate:"required,email"`
	Address   string `json:"address"  validate:"omitempty,min=6"`
	Password  string `gorm:"not null" json:"password" validate:"required,min=8,max=32"`
	Role      string `gorm:"default:admin"`
}
type OtpRecord struct {
	gorm.Model
	Value      string `json:"otp" validate:"required" `
	Email      string `json:"email" validate:"required,email"`
	ExpiryTime  time.Time
	Varify     bool   `gorm:"default:false"`
	Purpose      string `gorm:"default:account creation"`
}