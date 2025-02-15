package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PlaceOrder(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid user.Please login again",
		})
		return
	}
	UserId, ok := user_id.(int)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid user id"})
		return
	}
	userId := uint(UserId)
	var orderReq models.OrderRequest
	if err := c.BindJSON(&orderReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
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
		Where("users_id=?", userId).Find(&cart)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empty cart.Add Product to the cart"})
		return
	}
	var TotalAmount float64
	var DiscountAmount float64
	var FinalAmount float64
	var TotalProductCount uint

	for _, items := range cart {

		total := items.Price * float64(items.Quantity)
		discount := items.Discount / 100 * items.Price * float64(items.Quantity)
		finalAmt := total - discount

		TotalProductCount += items.Quantity
		TotalAmount += total
		DiscountAmount += discount
		FinalAmount += finalAmt
	}

	if orderReq.PaymentMethod == "Cash on delivery" && TotalAmount > 5000 {
		c.JSON(400, gin.H{"error": "Maximum value for the Cash on deliver 5000"})
		return
	}

	var ShippingAdd models.ShippingAddress
	result = database.DB.Model(&ShippingAdd).Where("user_id", userId).Find(&ShippingAdd)
	if result.RowsAffected < 1 {
		fmt.Println("user address already updated")
	}
	ShippingAdd = models.ShippingAddress{
		UserID:         userId,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	order := models.Order{
		UserID:            userId,
		ShippingAddressID: ShippingAdd.ID,
		SubTotal:          TotalAmount,
		ProductQuantity:   TotalProductCount,
		TotalDiscount:     DiscountAmount,
		TotalAmount:       FinalAmount,
		OrderStatus:       models.Pending,
		PaymentMethod:     orderReq.PaymentMethod,
		PaymentStatus:     models.PaymentPending,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(500, gin.H{"error": "database error in creating order"})
		return
	}

	var orderItems models.OrderItem

	for _, items := range cart {
		orderItems = models.OrderItem{
			OrderID:   order.ID,
			ProductID: items.ProductID,
			Quantity:  int(items.Quantity),
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
	var deleteCart models.CartItems
	if err := database.DB.Where("users_id=?", userId).Delete(&deleteCart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "datatbase error in deleting cart items"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"code":    "StatusCreated(201)",
		"details": "New order is created",
		"Data": gin.H{
			"User":             userId,
			"Shipping Address": ShippingAdd,
			"OrderDetails":     order,
		},
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
	fmt.Println(OrderDetails)
	orderResponce := models.OrderResponse{
		OrderID:       OrderDetails.ID,
		OrderStatus:   OrderDetails.OrderStatus,
		OrderTime:     OrderDetails.CreatedAt,
		TotalQuantity: int(OrderDetails.ProductQuantity),
		TotalAmount:   OrderDetails.TotalAmount,
		PaymentMethod: OrderDetails.PaymentMethod,
	}

	var orderItemResponce []models.OrderItemResponce
	if err := database.DB.Model(&models.OrderItem{}).Select("order_items.product_id,order_items.ID,products.product_name,products.description,products.image_url,products.discount,products.price,products.size,order_items.quantity").
		Joins("JOIN products on order_items.product_id=products.id").
		Where("order_id=?", order_id).Find(&orderItemResponce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"deatails": "order details fetched",
		"data": gin.H{
			"user":          user,
			"order":         orderResponce,
			"ordered_items": orderItemResponce,
		},
	})

}

func ListOrders(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missig.Please login again"})
	}
	var myOrders []models.OrderResponse
	result := database.DB.Model(&models.Order{}).Where("user_id=?", userId).Find(&myOrders)
	if result.Error != nil {
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
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is not exist.Please login again "})
		c.Abort()
		return

	}
	orderID := c.Param("order_id")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
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
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk",
		"details": "user order  ancelled",
	})

}
