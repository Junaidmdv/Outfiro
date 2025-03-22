package controllers

import (
	"errors"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"
    "fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OrdersList(c *gin.Context) {
	var orders []models.OrderList

	if err := database.DB.Model(&models.Order{}).
		Select("orders.user_id,users.email,orders.id,orders.order_status,orders.product_quantity,orders.created_at,orders.total_amount,payments.payment_method,payments.payment_status").
		Joins("JOIN payments ON orders.id=payments.order_id").
		Joins("JOIN users ON orders.user_id=users.id").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "orders data fetched",
		"data":    orders,
	})

}

func ChangeOrderStaus(c *gin.Context) {
	order_id := c.Param("id")
	var order models.Order
	result := database.DB.Where("id=?", order_id).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid order id",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}
	}
	var orderStatus models.OrderStatusUpdate
	if err := c.BindJSON(&orderStatus); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if order.OrderStatus == models.Cancelled {
		c.JSON(400, gin.H{"error": "Order cancelled.Status cann't changed"})
		return
	}
	if order.OrderStatus == models.Return {
		c.JSON(400, gin.H{"error": "Product are returned.Status cann't be changed"})
		return
	}
	if err := utils.ValidationOrderStatus(order.OrderStatus, orderStatus.Status); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if orderStatus.Status == models.Cancelled {
		var order []models.OrderItem
		result = database.DB.Model(&order).Where("order_id", order_id).Find(&order)
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
	}
	order.OrderStatus = orderStatus.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"details": "product status changed",
	})

}

func OrderItemList(c *gin.Context) {
	order_Id := c.Param("id")
	orderID, err := strconv.Atoi(order_Id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid order id ",
			"errors": err.Error()})
		return
	}
	var orderedItems []models.OrderItemResponce

	result := database.DB.Model(&models.OrderItem{}).Select("order_items.product_id,order_items.ID,products.product_name,products.description,products.image_url,products.discount,products.price,products.size,order_items.quantity").
		Joins("JOIN products on order_items.product_id=products.id").
		Where("order_id=?", orderID).Find(&orderedItems)

	if result.Error != nil {
		if result.Error == gorm.ErrCheckConstraintViolated {
			c.JSON(403, gin.H{"error": "Invalid order Id"})
		} else if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "Empty orders"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	} 
	fmt.Println("Add some idea")
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "order item retrived",
		"data":    orderedItems,
	})
}
