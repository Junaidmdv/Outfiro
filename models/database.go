package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	Pending           string  = "PENDING"
	Processing        string  = "PROCESSING"
	Shipped           string  = "SHIPPED"
	Cancelled         string  = "CANCELLED"
	Delivered         string  = "DELIVERED"
	Return            string  = "RETURN"
	PaymentPending    string  = "Payment pending"
	PaymentPaid       string  = "Payment paid"
	PaymentFailed     string  = "Payment Failed"
	PaymentRefunded   string  = "Payment Refunded"
	MaxQuantity       int     = 10
	Cash_on_delivery  string  = "Cash on delivery"
	Online_payment    string  = "Online payment"
	Wallete_payment   string  = "Wallete Payment"
	ReferralPoints    float64 = 100
	MinReffaralPoint  float64 = 100
	ReferalPointsUnit float64 = 10

	RefferalOffer string  = "Refferal cash bonus"
	CashCredited  string  = "Amount Credited"
	CashDebited   string  = "Amount Debited"
	MaxiLimitCOD  float64 = 20000

	DeliveryCharge float64=50
)

type Categories struct {
	gorm.Model
	CategoryName string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	// DiscountOffer float64    `json:"categories_offer" gorm:"default=0"`
	Products []Products `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type Products struct {
	gorm.Model
	ProductName   string          `json:"name" validate:"required,min=2,max=100"`
	Description   string          `json:"description" validate:"required,min=3,max=100,required"`
	ImageUrl      string          `json:"image_url"  validate:"required,url"`
	Discount      float64         `json:"discount_offer" gorm:"default=0"`
	CategoryId    uint            `json:"categoryid" binding:"required"`
	Price         float64         `json:"price" binding:"required" validate:"numeric,gt=0,required "`
	Size          string          `json:"size" binding:"required"`
	StockQuantity uint            `json:"StockQuantity" validate:"required,numeric,gt=0"`
	HasOffer      bool            `gorm:"default:false"`
	OrderItem     []OrderItem     `gorm:"foreignKey:ProductID"`
	CartItem      []CartItems     `gorm:"foreignKey:ProductID"`
	Wishlist_item []Wishlist_item `gorm:"foreignKey:ProductID"`
	Review        []Review        `gorm:"foreignKey:ProductID"`
}

type Users struct {
	gorm.Model
	FirstName       string `gorm:"not null" json:"firstname" validate:"required,alpha_space"`
	LastName        string `gorm:"not null" json:"lastname"  validate:"required,alpha_space"`
	Email           string `gorm:"not null;unique_index" json:"email" validate:"email"`
	PhoneNumber     string `gorm:"not null" json:"phone_number" validate:"phone_number"`
	Password        string `json:"password"`
	Role            string `gorm:"default:user" json:"role"`
	IsBlocked       bool   `gorm:"default:false" json:"isblocked"`
	LoginType       string `gorm:"default:normal"`
	ReferralCode    string `gorm:"referral_code" binding:"omitempty"`
	ReferralPoint   int    `gorm:"default:0"`
	ReferredCode    string `gorm:"default:0"`
	WallteAmount    float64
	Orders          []Order         `gorm:"foreignKey:UserID"`
	Address         []Address       `gorm:"foreignKey:UsersID"`
	CartItems       []CartItems     `gorm:"foreignKey:UsersID"`
	Wishlist_item   []Wishlist_item `gorm:"foreignKey:UsersID"`
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
	SubTotal          float64 `json:"Total_product_amount" validate:"numeric"`
	TotalDiscount     float64 `json:"Total_discount"`
	CouponOffer       float64
	DeliveryCharge    float64
	TotalAmount       float64 `json:"Total_Amount"`
	CouponID          *uint   `gorm:"default:null" binding:"required"`
	// PaymentMethod     string      `json:"payment_method"`
	// PaymentStatus     string      `json:"payment_status"`
	OrderStatus string      `gorm:"order_status" json:"order_status"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	Payment     Payment     `gorm:"foreignKey:OrderID"`
	Coupon      Coupon      `gorm:"foreignKey:CouponID"`
}

type Payment struct {
	gorm.Model
	OrderID           uint
	PaymentStatus     string `json:"payment_status" validate:"oneof='PAID''PENDING'"`
	RazorPayOrderID   string
	RazorpayPaymentID string
	RazorPaySignature string
	PaymentMethod     string  `json:"payment_method" validate:"oneof='CASH ON DELIVERY' 'Online Payment' 'Wallete_payment'"`
	Amount            float64 `json:"amount"`
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
	OrderID         uint `gorm:"not null"`
	ProductID       uint `gorm:"not null"`
	OrderItemStatus string
	Quantity        int      `gorm:"not null"`
	Product         Products `gorm:"foreignKey:ProductID"`
}

type Wishlist_item struct {
	gorm.Model
	UsersID   uint
	ProductID uint
}

type Coupon struct {
	ID              uint
	CouponCode      string      `json:"coupan_code"`
	DiscountPercent float64     `json:"discount_rate" validate:"numeric,gte=1"`
	MinPurchase     float64     `json:"min_purchase" validate:"numeric,gte=1"`
	ExpiredAt       time.Time   `json:"expired_at"`
	Limit           uint        `json:"limit"`
	CouponLimit     CouponLimit `json:"foreignKey:CouponID"`
}

type WalleteHistory struct {
	UserID         uint
	Amount         float64
	Reason         string
	TransationType string
	gorm.Model
}

type RazorpayResponce struct {
	OrderID  string  `json:"order_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	KeyID    string  `json:"key_id"`
}
type PaymentVerificationRequest struct {
	RazorPayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorPayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}

type SalesReport struct {
	UserId          uint      `json:"user_id" gorm:"column:user_id"`
	OrderId         uint      `json:"order_id" gorm:"column:id"`
	Email           string    `json:"email" gorm:"column:email"`
	OrderDate       time.Time `json:"ordered_date" gorm:"column:created_at"`
	ProductQuantity uint      `json:"product_quantity" gorm:"column:product_quantity"`
	Discount        float64   `json:"discount_offer" gorm:"total_discount"`
	CouponOffer     float64   `json:"coupon_discount" gorm:"coupon_offer"`
	TotalAmount     float64   `json:"total_amount" gorm:"column:total_amount"`
	PaymentMethod   string    `json:"payment_method" gorm:"column:payment_method"`
	PaymentStatus   string    `json:"payment_status" gorm:"column:payment_status"`
	OrderStatus     string    `json:"order_status" gorm:"column:order_status"`
}

type SalesReportPDf struct {
	OrderDate     string
	Email         string
	Quantity      string
	Discount      string
	TotalAmount   string
	PaymentMethod string
}

type OrderInvoice struct {
	Item         string
	Descritption string
	Quantity     string
	Price        string
	DiscountRate string
	Total        string
}

type CouponLimit struct {
	UserID     uint
	CouponID   uint
	CouponUsed uint `gorm:"default:0"`
}

type Transaction struct {
	ID     uint
	UserID uint
	Type   string
	Source string
	Amount float64
	gorm.Model
}

type Review struct {
	ID        uint
	UserID    uint
	ProductID uint
	Rating    float64 
	Comments  string 
	gorm.Model
	User       Users `gorm:"foreignKey:UserID"`
}

