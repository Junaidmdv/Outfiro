package models

import (
	"gorm.io/gorm"
)

type UpdateCategoryRequest struct {
	CategoriesName *string `json:"name"`
	Description    *string `json:"description"`
}
type ProductRequest struct {
	ProductName    string  `json:"name" validate:"required,alpha_space,max=100,min=2"`  
	ImageUrl       string   `json:"image_url" validate:"required,url"`
	Description    string  `json:"description" validate:"required"`
	CategoriesName string  `json:"category_name" binding:"required"`
	Price          float64 `json:"price" validate:"required,numeric,gt=0"`
	Status         string  `json:"status" validate:"oneof='In stock''Out of stock',required"`
	Size           string  `json:"size" validate:"required,oneof='M''XL''L''S'"`
	Quntity        uint    `json:"quantity" validate:"numeric,gt=0"`
}
type AdminLogin struct {
	gorm.Model
	Email    string `json:"email" validate:"email,min=8,required"`
	Password string `json:"password" validate:"min=8,required"`
}

type SignuPlayload struct {
	gorm.Model
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
type GetUser struct {
	Id         uint
	FirstName  string
	LastName   string
	Email      string
	Created_at string
	Updated_at string
	IsBlocked  bool
}

type UpdateUser struct {
	FirstName string `json:"firstname" validate:"alpha"`
	LastName  string `json:"last_name" validate:"alpha"`
	Email     string `json:"email" validate:"email"`
}

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type UpadatProduct struct {
	ProductName string  `json:"name" validate:"omitempty,min=5,max=50,alpha_space"`
	Price       float32 `json:"price" validate:"omitempty,numeric,gt=0"`
	Status      string  `json:"status" validate:"omitempty,oneof='In stock''Out of stock"`
	Quantity    uint    `json:"quantity" validate:"omitempty,numeric,gt=0"`
}
type ForgotPassword struct{
	Email string `json:"email" validate:"email"`
}

type ResetPassword struct {
	Email           string `json:"email" validate:"email"`
	Password        string `json:"password" validate:"min=8,max=32,alphanum"`
	ConfirmPassword string `json:"confirm_password" validate:"alphanum,eqfield=Password"`
}