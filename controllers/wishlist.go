package controllers

import (
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddWishlistItems(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "user id is not exist"})
		return
	}
	UserID := userID.(int)
	product_id := c.Param("product_id")
	productID, err := strconv.Atoi(product_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product_id",
		})
		return
	}
	if err := database.DB.
		Model(models.Products{}).
		Select("COUNT(*)").
		Where("id=?", productID).
		Scan(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not exist"})
		return
	}
	var wishlistexist bool
	if err := database.DB.
		Model(models.Wishlist_item{}).
		Select("COUNT(*)").
		Where("product_id=? AND users_id=?", productID, userID).
		Scan(&wishlistexist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if wishlistexist {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product already exist"})
		return
	}

	newItem := models.Wishlist_item{
		UsersID:   uint(UserID),
		ProductID: uint(productID),
	}
	if err := database.DB.Create(&newItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed  to add wishlist item"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   "success",
		"deatails": "product is added to the Wishlist",
	})
}

func GetWishlist(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var email string
	database.DB.Model(&models.Users{}).Where("id", userId).Pluck("email", &email)

	var items []models.WishlistResponce
	//
	if err := database.DB.Model(&models.Wishlist_item{}).
    Select("wishlist_items.id, wishlist_items.product_id, products.product_name, products.price, products.stock_quantity, products.image_url").
    Joins("JOIN products ON wishlist_items.product_id = products.id").
    Find(&items).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
    return
}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   "statusok",
		"data": gin.H{
			"user_id":        userId,
			"user_email":     email,
			"Wishlist_items": items,
		},
	})

}

func DeleteWishilistItems(c *gin.Context) {
	Wishlist_id := c.Param("id")
	WishlistID, err := strconv.Atoi(Wishlist_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item _id"})
		return
	}
	var DltItem models.Wishlist_item
	if err := database.DB.Model(&models.Wishlist_item{}).Where("id=?", WishlistID).Delete("id", DltItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "statusok",
		"details": "wishlist item deleted",
	})
}
