package models

import (
	"time"

	"gorm.io/gorm"
)

type UpdateCategoryRequest struct {
	CategoriesName *string `json:"name"`
	Description    *string `json:"description"`
}
type ProductRequest struct {
	ProductName    string  `json:"name" validate:"required,alpha_space,max=100,min=2"`
	ImageUrl       string  `json:"image_url" validate:"required,url"`
	Description    string  `json:"description" validate:"required"`
	Discount       int     `json:"discount" validate:"gte=1,lte=100"`
	CategoriesName string  `json:"category_name" binding:"required"`
	Price          float64 `json:"price" validate:"required,numeric,gt=0"`
	Size           string  `json:"size" validate:"required,oneof='M''XL''L''S'"`
	StockQuantity  uint    `json:"stock_quantity" validate:"numeric,gt=0"`
}
type AdminLogin struct {
	gorm.Model
	Email    string `json:"email" validate:"email,min=8,required"`
	Password string `json:"password" validate:"password"`
}

type SignuPlayload struct {
	gorm.Model
	FirstName       string `json:"name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"email,required"`
	PhoneNumber     string `json:"phone_number" validate:"phone_number"`
	Password        string `json:"password" validate:"password"`
	ConfirmPassword string `json:"confirm password" validate:"required"`
}

type LoginPlayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type GetUser struct {
	Id          uint
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Created_at  string
	Updated_at  string
	IsBlocked   bool
}

type UpdateUser struct {
	FirstName   string `json:"firstname" validate:"alpha"`
	LastName    string `json:"last_name" validate:"alpha"`
	Email       string `json:"email" validate:"email"`
	PhoneNumber string `json:"phone_number" validate:"phone_number"`
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
	ProductName   string  `json:"name" validate:"omitempty,min=5,max=50,alpha_space"`
	Price         float32 `json:"price" validate:"omitempty,numeric,gt=0"`
	StockQuantity uint    `json:"stock_quantity" validate:"omitempty,numeric,gt=0"`
	Discount      int     `json:"discount" validate:"omitempty,gte=1,lte=100"`
}
type ForgotPassword struct {
	Email string `json:"email" validate:"email"`
}

type ResetPassword struct {
	Email           string `json:"email" validate:"email"`
	Password        string `json:"password" validate:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type ProductResponce struct {
	ID          uint `gorm:"id"`
	ProductName string
	Description string
	Discount    int
	ImageUrl    string
	CategoryId  uint
	Price       float64
	Size        string
	QuantityAvailable uint `json:"quantity_available" gorm:"column:stock_quantity"`
}

type ProfileUpdate struct {
	ID          uint   `gorm:"id"`
	FirstName   string `json:"first_name"   validate:"omitempty,alpha_space"`
	LastName    string `json:"last_name"    validate:"omitempty,alpha_space"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone_number"`
}

type ProfileResponce struct {
	ID          uint `gorm:"id"`
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"current_password"`
	NewPassword        string `json:"new_password" validate:"password"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"eqfield=NewPassword"`
}

type AddressResponce struct {
	ID             uint   `gorm:"id"`
	Address_line_1 string `json:"address_line_1" validate:"omitempty,min=5"`
	Address_line_2 string `json:"address_line_2" validate:"omitempty"`
	City           string `json:"city" validate:"omitempty,alpha_space"`
	Pincode        string `json:"pincode" validate:"omitempty,pincode"`
	Landmark       string `json:"landmark" validate:"omitempty"`
	Contact_number string `json:"contact_number"  validate:"omitempty,numeric"`
	State          string `json:"state" validate:"omitempty,alpha_space"`
	Country        string `json:"country" validate:"omitempty,alpha_space"`
}
type CartResponce struct {
	ID            uint `gorm:"column:id" json:"id"`
	ProductID     uint `gorm:"product_id"`
	ProductName   string
	Description   string
	ImageUrl      string
	Discount      float64
	Price         float64
	Size          string
	StockQuantity uint
	Quantity      int
}

type OrderRequest struct {
	AddressId     uint   `json:"address_id"`
	PaymentMethod string `json:"payment_method"`
}
type CartRequest struct {
	Quantity int `json:"quantity" binding:"numeric"`
}

type OrderResponse struct {
	OrderID       uint      `json:"order_id" gorm:"column:id"`
	OrderStatus   string    `json:"order_status" gorm:"column:order_status"`
	TotalQuantity int        `json:"products_quantity" gorm:"column:product_quantity"`
	OrderTime     time.Time `json:"order_time" gorm:"column:created_at"`
	TotalAmount   float64   `json:"total_amount" gorm:"column:total_amount"`
	PaymentMethod string    `json:"payment_method" gorm:"column:payment_method"`
}
type OrderItemResponce struct {
	ProductID   uint    `gorm:"order_items.product_id"`
	ID          uint    `gorm:"order_items.id"`
	ProductName string  `gorm:"products.product_name"`
	Description string  `gorm:"products.description"`
	ImageUrl    string  `gorm:"products.image_url"`
	Discount    int     `gorm:"products.discount"`
	Price       float64 `gorm:"products.price"`
	Size        string  `gorm:"products.size"`
	Quantity    int     `gorm:"order_items.quantity"`
}

type CartValidation struct {
	ProductID uint    `gorm:"product_id"`
	Quantity  uint    `gorm:"cart.quantity"`
	Price     float64 `gorm:"product.price"`
	Discount  float64 `gorm:"product.discount"`
}

type OrderList struct {
	UserId          uint    `json:"user_id" gorm:"column:user_id"`
	OrderId         uint    `json:"order_id" gorm:"column:id"`
	Email           string  `json:"email" gorm:"column:email"`
	OrderDate       string  `json:"ordered_date" gorm:"column:created_at"`
	ProductQuantity uint    `json:"product_quantity" gorm:"column:product_quantity"`
	TotalAmount     float64 `json:"total_amount" gorm:"column:total_amount"`
	PaymentMethod   string  `json:"payment_method" gorm:"column:payment_method"`
	PaymentStatus   string  `json:"payment_status" gorm:"column:payment_status"`
	OrderStatus     string  `json:"order_status" gorm:"column:order_status"`
}

type OrderStatusUpdate struct {
	Status string `json:"order_status"`
}
