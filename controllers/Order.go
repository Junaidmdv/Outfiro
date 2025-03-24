package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// var OID uint

func PlaceOrder(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	UserId, ok := user_id.(int)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(UserId)
	var orderReq models.OrderRequest
	if err := c.BindJSON(&orderReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println(orderReq.PaymentMethod)
	fmt.Println(orderReq.CouponCode)

	var payment models.Payment
	switch orderReq.PaymentMethod {
	case models.Cash_on_delivery:
		payment.PaymentMethod = orderReq.PaymentMethod
	case models.Online_payment:
		payment.PaymentMethod = orderReq.PaymentMethod
	case models.Wallete_payment:
		payment.PaymentMethod = orderReq.PaymentMethod
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	var address models.Address
	result := database.DB.Model(&address).Where("users_id=?", orderReq.AddressId).First(&address)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user address not found"})
		return
	}

	var cart []models.CartValidation
	res := database.DB.Model(&models.CartItems{}).Select("cart_items.product_id,cart_items.quantity,products.price,products.discount").
		Joins("JOIN products ON cart_items.product_id=products.id").
		Where("users_id=?", userID).Find(&cart)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empty cart.Add Product to the cart"})
		return
	}

	//calculation of the total amount
	var TotalAmount float64
	var DiscountAmount float64
	var FinalAmount float64
	var TotalProductCount uint
	var couponDiscount float64

	for _, items := range cart {
		total := items.Price * float64(items.Quantity)
		discount := items.Discount / 100 * items.Price * float64(items.Quantity)
		finalAmt := total - discount

		TotalProductCount += items.Quantity
		TotalAmount += total
		DiscountAmount += discount
		FinalAmount += finalAmt
	}
	var order models.Order
	var coupon models.Coupon
	var tempCouponID *uint
	//coupan validation
	if orderReq.CouponCode != "" {
		result := database.DB.Where("coupon_code=?", orderReq.CouponCode).First(&coupon)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid coupon code"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
		}
		if coupon.ExpiredAt.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "coupon expired"})
			return
		}
		if coupon.MinPurchase > TotalAmount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Coupan not reached the minimum limit"})
			return
		}
		var limit models.CouponLimit
		result = database.DB.Where("user_id=? AND coupon_id=?", userID, coupon.ID).First(&limit)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				limit.UserID = userID
				limit.CouponID = coupon.ID
				database.DB.Create(&limit)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
		}

		if limit.CouponUsed == coupon.Limit {
			c.JSON(http.StatusForbidden, gin.H{"error": "Coupon usage exceed"})
			return
		}
		if err := database.DB.Model(&models.CouponLimit{}).Where("user_id=? AND coupon_id=?", userID, coupon.ID).Update("coupon_used", gorm.Expr("coupon_used+?", 1)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		fmt.Println(coupon)
		couponDiscount = float64(coupon.DiscountPercent) / 100 * FinalAmount
		FinalAmount -= couponDiscount
		tempCouponID = &coupon.ID

	}

	//Adding shippping address
	var ShippingAdd models.ShippingAddress
	result = database.DB.Model(&ShippingAdd).Where("user_id", userID).Find(&ShippingAdd)
	if result.RowsAffected < 1 {
		fmt.Println("user address already updated")
	}
	ShippingAdd = models.ShippingAddress{
		UserID:         userID,
		Address_line_1: address.Address_line_1,
		Address_line_2: address.Address_line_2,
		City:           address.City,
		Pincode:        address.Pincode,
		Landmark:       address.Landmark,
		Contact_number: address.Contact_number,
		State:          address.State,
		Country:        address.Country,
	}
	if err := database.DB.Model(&ShippingAdd).Save(&ShippingAdd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error in shipping address"})
		return
	}

	//adding new order
	order = models.Order{
		UserID:            userID,
		ShippingAddressID: ShippingAdd.ID,
		SubTotal:          TotalAmount,
		ProductQuantity:   TotalProductCount,
		TotalDiscount:     DiscountAmount,
		CouponID:          tempCouponID,
		CouponOffer:       couponDiscount,
		DeliveryCharge:    models.DeliveryCharge,
		TotalAmount:       FinalAmount + models.DeliveryCharge,
		OrderStatus:       models.Pending,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(500, gin.H{"error": "database error in creating order"})
		return
	}

	//
	orderRes := models.OrderResponse{
		OrderID:        order.ID,
		OrderStatus:    order.OrderStatus,
		TotalQuantity:  int(order.ProductQuantity),
		OrderTime:      order.CreatedAt,
		DeliveryCharge: order.DeliveryCharge,
		TotalAmount:    order.TotalAmount,
	}

	// create new payment record

	if orderReq.PaymentMethod == models.Cash_on_delivery {
		if TotalAmount > models.MaxiLimitCOD {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Maximum value for the Cash on deliver %f", models.MaxiLimitCOD)})
			return
		} else {
			payment.PaymentMethod = models.Cash_on_delivery
			payment.PaymentStatus = models.PaymentPending
		}
	}
	//razorpay
	if orderReq.PaymentMethod == models.Online_payment {
		payment.PaymentMethod = models.Online_payment
		payment.PaymentStatus = models.PaymentPending
		// OID = order.ID
	}

	//Wallete payment//
	if orderReq.PaymentMethod == models.Wallete_payment {
		var walleteAmount float64
		if err := database.DB.Model(&models.Users{}).Where("id=?", userID).Pluck("wallte_amount", &walleteAmount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error fetch wallete amount"})
			return
		}
		walleteAmount = walleteAmount - TotalAmount
		if walleteAmount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient wallete amount"})
			return
		}

		if err := database.DB.Model(&models.Users{}).Where("id=?", userID).Update("wallte_amount", walleteAmount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the wallete amount"})
			return
		}

		wallete := models.WalleteHistory{
			UserID:         userID,
			Amount:         TotalAmount,
			Reason:         "Order Placed",
			TransationType: models.CashDebited,
		}
		database.DB.Create(&wallete)

		// err := AddTransationDetails(userID, models.CashDebited, models.Wallete_payment, TotalAmount)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed add transaction details in database"})
		// 	return
		// }

		payment.PaymentMethod = models.Wallete_payment
		payment.PaymentStatus = models.PaymentPaid

	}

	payment.OrderID = order.ID
	payment.Amount = order.TotalAmount
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the payment details"})
		return
	}

	//adding the order items in database
	var orderItems models.OrderItem
	for _, items := range cart {
		orderItems = models.OrderItem{
			OrderID:         order.ID,
			ProductID:       items.ProductID,
			Quantity:        int(items.Quantity),
			OrderItemStatus: order.OrderStatus,
		}
		database.DB.Create(&orderItems)
		var product models.Products
		database.DB.Model(&product).Where("id=?", items.ProductID).First(&product)
		product.StockQuantity = product.StockQuantity - items.Quantity
		database.DB.Save(&product)
	}
	var count int64
	database.DB.Model(&orderItems).Where("order_id=?", order.ID).Count(&count)
	if count == 0 {
		c.JSON(500, gin.H{"error": "No records are added"})
		return
	}
	//delete the cartitems
	var deleteCart models.CartItems
	if err := database.DB.Where("users_id=?", userID).Delete(&deleteCart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "datatbase error in deleting cart items"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"code":    "StatusCreated(201)",
		"details": "New order is created",
		"Data":    orderRes,
	})

}

func GetOrderDetails(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User id is missing.Login again"})
		return
	}
	var OrderDetails models.Order
	order_id := c.Param("order_id")
	OrderID, err := strconv.Atoi(order_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order id"})
		return
	}
	var user models.ProfileResponce
	if err := database.DB.Model(&models.Users{}).Where("id=?", user_id).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	result := database.DB.Model(&models.Order{}).Where("id=?", OrderID).First(&OrderDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order is not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error in order id",
			})
			return
		}
	}

	orderResponce := models.OrderResponse{
		OrderID:        OrderDetails.ID,
		OrderStatus:    OrderDetails.OrderStatus,
		OrderTime:      OrderDetails.CreatedAt,
		TotalQuantity:  int(OrderDetails.ProductQuantity),
		TotalAmount:    OrderDetails.TotalAmount,
		DeliveryCharge: OrderDetails.DeliveryCharge,
	}

	var orderItemResponce []models.OrderItemResponce
	if err := database.DB.Model(&models.OrderItem{}).Select("order_items.product_id,order_items.ID,products.product_name,products.description,products.image_url,products.discount,products.price,products.size,order_items.quantity,order_items.order_item_status").
		Joins("JOIN products on order_items.product_id=products.id").
		Where("order_id=?", order_id).Find(&orderItemResponce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var shippingAddress models.ShoppingAddressResponce
	if err := database.DB.Model(&models.ShippingAddress{}).Where("id=?", OrderDetails.ShippingAddressID).First(&shippingAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "datatbase  error"})
		return
	}

	var order models.Order
	if err := database.DB.Model(&order).Where("id=?", OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	var payment models.PaymentResponce
	if err := database.DB.Model(&models.Payment{}).Where("order_id=?", OrderID).First(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get payment details"})
		return
	}

	if order.CouponOffer > 0 {
		c.JSON(http.StatusOK, gin.H{
			"Status":  "success",
			"details": "order details fetched",
			"data": gin.H{
				"user":            user,
				"order details":   orderResponce,
				"coupon":          "applied",
				"coupon discount": order.CouponOffer,
				"payment_method":  payment.PaymentMethod,
				"payment_status":  payment.PaymentStatus,
				"ordered items":   orderItemResponce,
				"shipping":        shippingAddress,
			},
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"deatails": "order details fetched",
		"data": gin.H{
			"user":           user,
			"order":          orderResponce,
			"payment_method": payment.PaymentMethod,
			"payment_status": payment.PaymentStatus,
			"ordered_items":  orderItemResponce,
			"shipping":       shippingAddress,
		},
	})

}

func ListOrders(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missig.Please login again"})
	}
	var myOrders []models.OrderResponse

	if err := database.DB.Model(&models.Order{}).
		Select("orders.id,orders.order_status,orders.product_quantity,orders.created_at,orders.total_amount,payments.payment_method,payments.payment_status").
		Joins("JOIN payments ON orders.id=payments.order_id").
		Where("user_id=?", userId).
		Find(&myOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "order details are fetched",
		"data":    myOrders,
	})
}
func CancelOrder(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is not exist.Please login again "})
		c.Abort()
		return
	}
	userID := userId.(int)
	orderID := c.Param("id")
	var CancelOrder models.Order
	result := database.DB.Model(&models.Order{}).Where("id", orderID).First(&CancelOrder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": "Orders not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetch order details"})
			return
		}
	}
	if CancelOrder.OrderStatus == models.Cancelled {
		c.JSON(400, gin.H{"error": "orders already cancelled"})
		return
	}

	if CancelOrder.OrderStatus == models.Delivered {
		c.JSON(400, gin.H{"error": "orders already delivered"})
		return
	}
	if CancelOrder.OrderStatus == models.Shipped {
		c.JSON(400, gin.H{"error": "orders shipped.Cann't cancel the orders"})
		return
	}
	CancelOrder.OrderStatus = models.Cancelled
	if err := database.DB.Save(&CancelOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var order []models.OrderItem
	result = database.DB.Model(&order).Where("order_id", orderID).Find(&order)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order items are empty"})
		return
	}
	for _, items := range order {
		var product models.Products
		if err := database.DB.First(&product, items.ProductID).Error; err != nil {
			c.JSON(400, gin.H{"error": "Product not found"})
			return
		}
		if err := database.DB.Model(&product).Update("stock_quantity", gorm.Expr("stock_quantity + ?", items.Quantity)).Error; err != nil {
			c.JSON(400, gin.H{"error": "Failed to update product quantity"})
			return
		}
		if err := database.DB.Model(&items).Where("id=?", items.ID).Update("order_item_status", CancelOrder.OrderStatus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error in updating user items status"})
			return
		}
	}
	var payment models.Payment
	if err := database.DB.Model(&models.Payment{}).Where("order_id=?", orderID).First(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}
	if payment.PaymentMethod == models.Wallete_payment || payment.PaymentMethod == models.Online_payment {
		//update payment when the wallete or cancel
		payment.PaymentStatus = models.PaymentRefunded
		payment.Amount = 0

		//add refunded money to the wallete
		if err := database.DB.Model(&models.Users{}).Where("id=?", userId).Update("wallte_amount", gorm.Expr("wallte_amount+?", CancelOrder.TotalAmount)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the refunded amount in the wallete"})
			return
		}
		//add record on  WalleteHistory
		RefundWallete := models.WalleteHistory{
			UserID:         uint(userID),
			Amount:         CancelOrder.TotalAmount,
			Reason:         models.PaymentRefunded,
			TransationType: models.CashCredited,
		}
		if err := database.DB.Create(&RefundWallete).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed add new wallete hisory record"})
			return
		}

	}
	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the payment"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk",
		"details": "user order  cancelled",
	})

}

func CancelOrderProduct(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	userID := user_id.(int)

	Order_item_id := c.Param("id")
	OrderitemID, err := strconv.Atoi(Order_item_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user item id"})
		return
	}

	var cancelOrderItem models.OrderItem
	result := database.DB.Where("id=?", OrderitemID).Preload("Product").First(&cancelOrderItem)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order item not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "StatusInternalServerError"})
			return
		}
	}

	if cancelOrderItem.OrderItemStatus == models.Shipped || cancelOrderItem.OrderItemStatus == models.Delivered || cancelOrderItem.OrderItemStatus == models.Cancelled || cancelOrderItem.OrderItemStatus == models.Return {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User cann't cancel the product.Product already %s", cancelOrderItem.OrderItemStatus)})
		return
	}
	cancelOrderItem.OrderItemStatus = models.Cancelled

	CancelOrderAmount := cancelOrderItem.Product.Price * float64(cancelOrderItem.Quantity)

	//update the order amount
	var order models.Order
	if err := database.DB.Model(&order).Where("id=?", cancelOrderItem.OrderID).Preload("Coupon").First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	order.ProductQuantity -= uint(cancelOrderItem.Quantity)
	if order.ProductQuantity == 0 {
		order.OrderStatus = models.Cancelled
	}

	if order.CouponID != nil {

		TotalAmount := order.SubTotal - CancelOrderAmount
		if TotalAmount < order.Coupon.MinPurchase {

			order.TotalAmount = TotalAmount + order.DeliveryCharge
			order.CouponOffer = 0
		} else {

			discountAmount := TotalAmount * order.Coupon.DiscountPercent / 100
			order.TotalAmount = TotalAmount - discountAmount + order.DeliveryCharge
			order.CouponOffer = discountAmount
		}

	}
	if order.CouponID == nil {
		order.TotalAmount = order.TotalAmount - CancelOrderAmount
	}
	order.SubTotal = order.SubTotal - CancelOrderAmount

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update the total amount"})
		return
	}

	if err := database.DB.Model(&models.Products{}).Where("id=?", cancelOrderItem.ProductID).Update("stock_quantity", gorm.Expr("stock_quantity+?", cancelOrderItem.Quantity)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the product stock"})
		return
	}

	var payment models.Payment

	database.DB.Where("order_id=?", cancelOrderItem.OrderID).First(&payment)

	if payment.PaymentMethod == models.Wallete_payment || payment.PaymentMethod == models.Online_payment {
		if err := database.DB.Model(&models.Users{}).Where("id=?", userID).Update("WallteAmount", gorm.Expr("wallte_amount+?", CancelOrderAmount)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update wallete amount"})
			return
		}
		Wallete := models.WalleteHistory{
			UserID:         uint(userID),
			Amount:         CancelOrderAmount,
			Reason:         "Order Product cancelled",
			TransationType: models.CashCredited,
		}
		if err := database.DB.Model(&models.WalleteHistory{}).Create(&Wallete).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
	if err := database.DB.Save(&cancelOrderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"details": "product canceled",
	})
}

func OrderReturn(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userId := userID.(int)
	orderId := c.Param("id")
	fmt.Println(orderId)
	orderID, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}
	var orderReturn models.Order
	result := database.DB.Model(&orderReturn).Where("id=? AND user_id=?", orderID, userID).First(&orderReturn)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order is not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error in order id",
			})
			return
		}
	}
	if orderReturn.OrderStatus != models.Delivered {
		c.JSON(http.StatusForbidden, gin.H{"error": "order status should be delivered to return the product"})
		return
	}

	if orderReturn.OrderStatus == models.Return {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order already returned"})
		return
	}
	var payment models.Payment
	database.DB.Where("order_id=?", orderID).First(&payment)
	payment.PaymentStatus = models.PaymentRefunded
	database.DB.Save(&payment)

	var user models.Users
	database.DB.Where("user_id=?", userID).Find(&user)
	user.WallteAmount = user.WallteAmount + payment.Amount
	database.DB.Save(&user)

	wallete := models.WalleteHistory{
		UserID:         uint(userId),
		Amount:         payment.Amount,
		Reason:         "Order Returned",
		TransationType: models.CashCredited,
	}
	if err := database.DB.Create(&wallete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

}

func WalleteHistory(c *gin.Context) {
	userId, _ := c.Get("user_id")
	// userID := userId.(int)

	var wallete []models.WalleteResponce
	if err := database.DB.Model(&models.WalleteHistory{}).Where("user_id=?", userId).Find(&wallete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	var email string
	database.DB.Model(&models.Users{}).Where("id=?", userId).Pluck("email", &email)

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"user email": email,
		"data":       wallete,
	})
}

// func AddTransationDetails(UserID uint, Type string, Source string, Amount float64) error {
// 	record := models.Transaction{
// 		UserID: UserID,
// 		Type:   Type,
// 		Source: Source,
// 		Amount: Amount,
// 	}
// 	if err := database.DB.Create(&record).Error; err != nil {
// 		return fmt.Errorf("failed add transaction details")
// 	}
// 	return nil
// }

// func TransactionHistory(c *gin.Context) {

// }
