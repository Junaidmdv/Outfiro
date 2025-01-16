package models

import (
	"gorm.io/gorm"
)

type Categories struct {
	gorm.Model
	CategoryName string     `json:"name" binding:"required"`
	Description  string     `json:"description"`
	ImageUrl     string     `json:"imageurl" binding:"required"`
	Products     []Products `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UpdateCategoryRequest struct {
	CategoriesName *string `json:"name"`
	Description    *string `json:"description"`
}

type Products struct {
	gorm.Model
	ProductName string  `json:"product name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	CategoryId  uint    `json:"categoryid" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Status      string  `type:"enum('in stock','out of stock');" json:"status"`
	Size        string  `json:"size" binding:"required"`
	Quntity     string  `json:"quantity" binding:"required"`
	Discount    string  `json:"discount,omitempty" `
}

type ProductRequest struct {
	ProductName    string  `json:"product name" binding:"required"`
	Description    string  `json:"description" binding:"required"`
	CategoriesName string  `json:"categoryname" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	Status         string  `type:"enum('in stock','out of stock');" json:"status"`
	Size           string  `json:"size" binding:"required"`
	Quntity        string  `json:"quantity" binding:"required"`
	Discount       string  `json:"discount,omitempty" `
}

type Users struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname" validate:"required,min=4,alpha"`
	LastName  string `gorm:"not null" json:"lastname"  validate:"required alpha"`
	Email     string `gorm:"not null;unique_index" json:"email" validate:"email"`
	Address   string `gorm:"not null"  json:"address" validate:"omitempty"`
	Password  string `gorm:"not null" json:"password"`
	Role      string `gorm:"not null,default:user" json:"role"`
	IsBlocked bool   `gorm:"default:false" json:"isblocked"`
}
type Admin struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname" validate:"required,alpha,max=6"`
	LastName  string `gorm:"not null" json:"lastname"  validate:"omitempty"`
	Email     string `gorm:"not null;unique_index" json:"email"  validate:"required,email"`
	Address   string `gorm:"not null"  json:"address"  validate:"omitempty,min=6"`
	Password  string `gorm:"not null" json:"password" validate:"required,min=8,max=32,ContainsSpecialChars"`
	Role      string `gorm:"default:admin"`
}
type AdminLogin struct {
	gorm.Model
	Email    string `json:"email" validate:"email,max=8,required"`
	Password string `json:"password" validate:"max=8,required"`
}

type SignuPlayload struct {
	FirstName       string `json:"name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm password" binding:"required"`
}

type LoginPlayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
