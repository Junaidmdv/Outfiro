package controllers

import (
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"

	"github.com/gin-gonic/gin"
)

func GetSalesData(c *gin.Context) {
	// filter := c.Query("filter")
	var TotalUser uint
	if err := database.DB.Model(&models.Users{}).Select("COUNT(id) as total_user").Pluck("total_user", &TotalUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	var order []models.Order
	if err := database.DB.Model(&models.Order{}).Preload("Payment").Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	var SalesAmount float64
	var TotalDiscount float64
	var CouponDiscount float64

	Orders := map[string]int{
		"orders":          0,
		"order_delivered": 0,
		"order_cancelled": 0,
	}

	for _, items := range order {
		SalesAmount += items.TotalAmount
		TotalDiscount += items.TotalDiscount
		CouponDiscount += items.CouponOffer
		Orders["orders"]++
		if items.OrderStatus == models.Cancelled {
			Orders["order_cancelled"]++
		}
		if items.OrderStatus == models.Delivered {
			Orders["order_delivered"]++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"Order":          Orders,
			"total_sales":    SalesAmount,
			"total_discount": TotalDiscount,
			"coupon_offers":  CouponDiscount,
		},
	})

}

func BestSellingProduct(c *gin.Context) {
	var BestSellingProduct []models.BestSellingProduct
	if err := database.DB.Table("order_items").Limit(10).
		Where("order_items.order_item_status <> ? AND order_items.deleted_at IS NULL", models.Cancelled).
		Select("order_items.product_id, products.product_name, products.description, products.image_url, products.price,SUM(order_items.quantity) as product_sold").
		Joins("JOIN products ON order_items.product_id = products.id").
		Group("order_items.product_id,products.product_name,products.description,products.image_url,products.price").
		Order("product_sold DESC").
		Find(&BestSellingProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   BestSellingProduct,
	})

}

func BestCategories(c *gin.Context) {
	var categories []models.BestSellingCategories
	if err := database.DB.Table("categories").Limit(10).
		Where("order_items.order_item_status <> ? AND order_items.deleted_at IS NULL", models.Cancelled).
		Select("products.category_id,categories.category_name,categories.description,SUM(order_items.quantity) total_product_sold").
		Joins("JOIN products ON categories.id=products.category_id").
		Joins("JOIN order_items ON products.id=order_items.product_id").
		Group("products.category_id,categories.category_name,categories.description").
		Order("total_product_sold").
		Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   categories,
	})

}

func SalesGraph(c *gin.Context) {
	filter := c.Query("filter")
	fmt.Println(filter)

	var SalesData []models.SaleGraphData
	switch filter {
	case "year":
		database.DB.Model(models.Order{}).Select("TO_CHAR(created_at,'YYYY') as date,SUM(total_amount) as sales").
			Group("date").
			Order("date").
			Find(&SalesData)
	case "month":
		database.DB.Model(models.Order{}).Select("TO_CHAR(created_at,'YYYY-MM') as date,SUM(total_amount) as sales").
			Group("date").
			Order("date").
			Find(&SalesData)

	case "week":
		database.DB.Model(models.Order{}).Select("TO_CHAR(created_at,'IYYY-IW') as date,SUM(total_amount) as sales").
			Group("date").
			Order("date").
			Find(&SalesData)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter"})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"Data":   SalesData,
	})

}
