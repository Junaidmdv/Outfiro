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
	ProductName    string  `json:"name" validate:"required,max=100,min=2"`
	ImageUrl       string  `json:"image_url" validate:"required,url"`
	Description    string  `json:"description" validate:"required"`
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
	ReferedCode     string `json:"referral_code" validate:"omitempty,gte=8"`
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
	ID                uint `gorm:"id"`
	ProductName       string
	Description       string
	Discount          int
	ImageUrl          string
	CategoryId        uint
	Price             float64
	Size              string
	QuantityAvailable uint `json:"quantity_available" gorm:"column:stock_quantity"`
}

type ProductListing struct {
	ID            uint `gorm:"id"`
	ProductName   string
	Discount      int
	ImageUrl      string
	Price         float64
	Size          string
	StockQuantity uint    `json:"inventory_details"`
	AvgRating     float64 `json:"avg_rating" gorm:"column:avg_ratings"`
}

type ProfileUpdate struct {
	ID          uint   `gorm:"id"`
	FirstName   string `json:"first_name"   validate:"omitempty,alpha_space"`
	LastName    string `json:"last_name"    validate:"omitempty,alpha_space"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone_number"`
}

type ProfileResponce struct {
	ID            uint `gorm:"id"`
	FirstName     string
	LastName      string
	PhoneNumber   string
	Email         string
	ReferralCode  string `gorm:"referral_code"`
	ReferralPoint uint
	WallteAmount  float64
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
	CouponCode    string `json:"coupon_code"`
}
type CartRequest struct {
	Quantity int `json:"quantity" binding:"numeric"`
}

type OrderResponse struct {
	OrderID        uint      `json:"order_id" gorm:"column:id"`
	OrderStatus    string    `json:"order_status" gorm:"column:order_status"`
	TotalQuantity  int       `json:"products_quantity" gorm:"column:product_quantity"`
	OrderTime      time.Time `json:"order_time" gorm:"column:created_at"`
	DeliveryCharge float64   `gorm:"delivery_charge" json:"delivery_charge"`
	TotalAmount    float64   `json:"total_amount" gorm:"column:total_amount"`
}
type OrderItemResponce struct {
	ProductID       uint    `gorm:"order_items.product_id"`
	ID              uint    `gorm:"order_items.id" json:"OrderItems_id"`
	ProductName     string  `gorm:"products.product_name"`
	Description     string  `gorm:"products.description"`
	ImageUrl        string  `gorm:"products.image_url"`
	Discount        int     `gorm:"products.discount"`
	Price           float64 `gorm:"products.price"`
	Size            string  `gorm:"products.size"`
	Quantity        int     `gorm:"order_items.quantity"`
	OrderItemStatus string  `gorm:"order_item_status" json:"product_status"`
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
type WishlistResponce struct {
	ID        uint   `gorm:"id"`
	ProductID uint   `gorm:"column:product_id" json:"product_id"`
	Name      string `gorm:"column:product_name" json:"product_name"`
	Price     string `gorm:"price" json:"price"`
	Quantity  int    `gorm:"column:stock_quantity" json:"stock quantity"`
	ImageUrl  string `gorm:"column:image_url"`
}

// type RazorpayPayment struct {
// 	OrderID     string
// 	Email       string
// 	PhoneNumber string
// 	TotalAmount float64
// }

type RazorpayClaims struct {
	OrderID   string `json:"razorpay_order_id"`
	PaymentID string `json:"razorpay_payment_id"`
	Signature string `json:"razorpay_signature"`
}

type ReferalPoints struct {
	Amount float64 `json:"referal_points"`
}
type CouponUpdate struct {
	DiscountPercent float64 `json:"discount_rate" validate:"numeric,gt=0"`
	MinPurchase     float64 `json:"min_purchase" validate:"numeric,gt=0"`
	UsageLimit      uint    `json:"maximum_limit" validate:"numeric"`
}

type SalesReportRequest struct {
	StartDate string `json:"start_date" binding:"omitempty"`
	EndDate   string `json:"end_date"  binding:"omitempty"`
	Limit     string `json:"limit"   binding:"omitempty"`
}

type WalleteResponce struct {
	Amount         float64
	Reason         string
	TransationType string
}

type ShoppingAddressResponce struct {
	UserID         uint
	Address_line_1 string `json:"address_line_1" validate:"required"`
	Address_line_2 string `json:"address_line_2" validate:"required"`
	City           string `json:"city" validate:"required,alpha_space"`
	Pincode        string `json:"pincode" validate:"required,pincode"`
	Landmark       string `json:"landmark" validate:"required"`
	Contact_number string `json:"contact_number"  validate:"required,numeric"`
	State          string `json:"state" validate:"required,alpha_space"`
	Country        string `json:"country" validate:"required,alpha_space"`
}

type OfferResponce struct {
	ProductName string  `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Description string  `json:"description" validate:"required,min=3,max=100,required"`
	ImageUrl    string  `json:"image_url"  validate:"required,url"`
	Discount    float64 `json:"discount_offer" gorm:"default=0"`
	CategoryId  uint    `json:"categoryid" binding:"required"`
	Price       float64 `json:"price" binding:"required" validate:"numeric,gt=0,required "`
	Size        string  `json:"size" binding:"required"`
}

type BestSellingProduct struct {
	ProductID   uint
	ProductName string
	ImageUrl    string
	Price       float64
	ProductSold uint
}

type BestSellingCategories struct {
	ID           uint   `json:"ID" gorm:"column:category_id"`
	CategoryName string `json:"name"`
	Description  string `json:"description"`
	ProductSold  uint   `json:"total_sales" gorm:"column:total_product_sold"`
}

type ReviewRequest struct {
	Rating   float64 `json:"rating" validate:"omitempty,gte=1,lte=5"`
	Comments string  `json:"comments" validate:"omitempty,min=2,max=100"`
}

type ReviewResponce struct {
	ID         uint
	UserName   string
	Rating     float64
	Comments   string
	Created_at string
}

type SaleGraphData struct {
	Date  string `gorm:"column:date"`
	Sales string `gorm:"column:sales"`
}

type PaymentResponce struct {
	PaymentStatus string `json:"payment_status"`
	PaymentMethod string `json:"payment_method"`
}
