package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	Pending        string = "PENDING"
	Processing     string = "PROCESSING"
	Shipped        string = "SHIPPED"
	Cancelled      string = "CANCELLED"
	Delivered      string = "DELIVERED"
	Return         string = "RETURN"
	PaymentPending string = "Payment pending"
	PaymentPaid    string = "Payment paid"
	MaxQuantity    int    = 10
)

type Categories struct {
	gorm.Model
	CategoryName string     `json:"name" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	Products     []Products `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type Products struct {
	gorm.Model
	ProductName   string    `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Description   string    `json:"description" validate:"required,min=3,max=100,required"`
	ImageUrl      string    `json:"image_url"  validate:"required,url"`
	Discount      int       `json:"discount" validate:"gte=0,lte=100"`
	CategoryId    uint      `json:"categoryid" binding:"required"`
	Price         float64   `json:"price" binding:"required" validate:"numeric,gt=0,required "`
	Size          string    `json:"size" binding:"required"`
	StockQuantity uint      `json:"StockQuantity" validate:"required,numeric,gt=0"`
	OrderItem     OrderItem `gorm:"foreignKey:ProductID"`
	CartItem      CartItems `gorm:"foreignKey:ProductID"`
}
type Users struct {
	gorm.Model
	FirstName       string          `gorm:"not null" json:"firstname" validate:"required,alpha_space"`
	LastName        string          `gorm:"not null" json:"lastname"  validate:"required,alpha_space"`
	Email           string          `gorm:"not null;unique_index" json:"email" validate:"email"`
	PhoneNumber     string          `gorm:"not null" json:"phone_number" validate:"phone_number"`
	Password        string          `json:"password"`
	Role            string          `gorm:"default:user" json:"role"`
	IsBlocked       bool            `gorm:"default:false" json:"isblocked"`
	LoginType       string          `gorm:"default:normal"`
	Orders          []Order         `gorm:"foreignKey:UserID"`
	Address         []Address       `gorm:"foreignKey:UsersID"`
	CartItems       []CartItems     `gorm:"foreignKey:UsersID"`
	ShippingAddress ShippingAddress `gorm:"foreignKey:UserID"`
}
type Address struct {
	gorm.Model
	UsersID        uint   `json:"user_id" validate:"required,numeric,gte=1"`
	Address_line_1 string `json:"address_line_1" validate:"required"`
	Address_line_2 string `json:"address_line_2" validate:"required"`
	City           string `json:"city" validate:"required,alpha_space"`
	Pincode        string `json:"pincode" validate:"required,pincode"`
	Landmark       string `json:"landmark" validate:"required"`
	Contact_number string `json:"contact_number"  validate:"required,numeric"`
	State          string `json:"state" validate:"required,alpha_space"`
	Country        string `json:"country" validate:"required,alpha_space"`
}

type Admin struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname" validate:"required"`
	LastName  string `gorm:"not null" json:"lastname"  validate:"omitempty"`
	Email     string `gorm:"not null;unique_index" json:"email"  validate:"required,email"`
	Address   string `json:"address"  validate:"omitempty,min=6"`
	Password  string `gorm:"not null" json:"password" validate:"password"`
	Role      string `gorm:"default:admin"`
}
type OtpRecord struct {
	gorm.Model
	Value      string `json:"otp" validate:"required" `
	Email      string `json:"email" validate:"required,email"`
	ExpiryTime time.Time
	Varify     bool   `gorm:"default:false"`
	Purpose    string `gorm:"default:account creation"`
}

type CartItems struct {
	gorm.Model
	UsersID   uint
	ProductID uint
	Quantity  int `json:"quantity" validate:"numeric,gte=0"`
}

type Order struct {
	gorm.Model
	UserID            uint
	ShippingAddressID uint
	ProductQuantity   uint
	SubTotal          float64     `json:"Total_product_amount" validate:"numeric"`
	TotalDiscount     float64     `json:"Total_discount"`
	TotalAmount       float64     `json:"Total_Amount"`
	PaymentMethod     string      `json:"payment_method"  validate:"oneof='Cash on deliver'"`
	PaymentStatus     string      `json:"payment_status"`
	OrderStatus       string      `gorm:"order_status" json:"order_status"`
	OrderItems        []OrderItem `gorm:"foreignKey:OrderID"`
}

type ShippingAddress struct {
	gorm.Model
	UserID         uint
	Address_line_1 string `json:"address_line_1" validate:"required"`
	Address_line_2 string `json:"address_line_2" validate:"required"`
	City           string `json:"city" validate:"required,alpha_space"`
	Pincode        string `json:"pincode" validate:"required,pincode"`
	Landmark       string `json:"landmark" validate:"required"`
	Contact_number string `json:"contact_number"  validate:"required,numeric"`
	State          string `json:"state" validate:"required,alpha_space"`
	Country        string `json:"country" validate:"required,alpha_space"`
	Order          Order  `gorm:"foreignKey:ShippingAddressID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
}
