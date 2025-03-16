package controllers

import (
	"errors"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProductOffter(c *gin.Context) {
	ProductId := c.Param("product_id")
	ProductID, err := strconv.Atoi(ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid produdct id"})
		return
	}
	var Product models.Products
	result := database.DB.Model(&Product).Where("id=?", ProductID).First(&Product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found Or "})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return
		}
	}
	if Product.HasOffer {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Offer already applied",
		})
		return
	}

	type ProductOfferdiscount struct {
		Value float64 `json:"discount_offer"`
	}
	var discountPercent ProductOfferdiscount

	if err := c.BindJSON(&discountPercent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update the json"})
		return
	}

	if discountPercent.Value < 1 && discountPercent.Value > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage"})
		return
	}
	Product.Discount = discountPercent.Value
	Product.HasOffer = true
	database.DB.Save(&Product)
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "Discount offer is added to the product",
	})
}

func ListOffer(c *gin.Context) {
	var offerProduct []models.OfferResponce
	if err := database.DB.Model(&models.Products{}).Where("has_offer", true).Find(&offerProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"Data":        offerProduct,
		"Total offer": len(offerProduct),
	})

}

func RemoveOffer(c *gin.Context) {
	ProductId := c.Param("product_id")
	ProductID, err := strconv.Atoi(ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid produdct id"})
		return
	}
	var Product models.Products
	result := database.DB.Model(&Product).Where("id=? AND has_offer=?", ProductID, true).First(&Product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not applied offer"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return
		}
	}
	if !Product.HasOffer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offer already removed Or not appleid"})
		return
	}

	Product.Discount = 0
	Product.HasOffer = false
	database.DB.Save(&Product)
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "Discount offer removed",
	})
}

// func AddCategorieOffer(c *gin.Context) {
// 	categoryId := c.Param("category_id")
// 	categoryID, err := strconv.Atoi(categoryId)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid categories id"})
// 		return
// 	}
// 	var category models.Categories
// 	result := database.DB.Model(&category).Where("id=?", categoryID)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
// 			return
// 		}
// 	}
// 	type CategoryOffer struct {
// 		Value float64
// 	}
// 	var categoryDiscount CategoryOffer
// 	if err := c.BindJSON(&categoryDiscount); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category offer"})
// 		return
// 	}
// 	if categoryDiscount.Value < 1 && categoryDiscount.Value > 100 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage"})
// 		return
// 	}
// 	category.DiscountOffer = categoryDiscount.Value
// 	database.DB.Save(&categoryDiscount)
// 	c.JSON(200, gin.H{
// 		"status":  "success",
// 		"details": "Discount offer is added to the product",
// 	})

// }

