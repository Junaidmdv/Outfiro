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

func AddToCart(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User id is invalid"})
		return
	}

	user_id, ok := userId.(int)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid User id"})
		return
	}

	product_id := c.Param("product_id")

	ProductId, err := strconv.Atoi(product_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id",
			"errors": err.Error()})
		return
	}
	var product models.Products

	result := database.DB.Where("id=?", ProductId).First(&product)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var cartReq models.CartRequest
	if err := c.BindJSON(&cartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if cartReq.Quantity <= 0 {
		c.JSON(400, gin.H{"error": "Quantity should be at least 1. The item will not be added to cart"})
		return
	}
	if cartReq.Quantity > models.MaxQuantity {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Product quantity exceed maximum limit.Only 10 product quantity availible"})
		return
	}
	if product.StockQuantity == 0 {
		c.JSON(400, gin.H{"error": "Product of stock"})
		return
	}
	quantityleft := int(product.StockQuantity) - cartReq.Quantity
	if quantityleft < 0 {
		c.JSON(400, gin.H{"error": fmt.Sprintf("The requested stock is not availible.Only %d stock is availible", product.StockQuantity)})
		return
	}
	var cartItem models.CartItems
	result = database.DB.Model(&cartItem).Where("product_id=?", product_id).First(&cartItem)
	if result.RowsAffected > 0 {
		c.JSON(400, gin.H{"error": "product already exist in the cart"})
		return
	}

	cartItem.Quantity = cartReq.Quantity
	cartItem.ProductID = product.ID
	cartItem.UsersID = uint(user_id)

	if err := database.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart items"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"code":    "StatusCreated",
		"details": "Product added to the cart",
	})

}

func CartRemoveItem(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missging"})
		return
	}
	ProductIDstr := c.Param("product_id")
	ProductID, err := strconv.Atoi(ProductIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user product id"})
	}
	var removeItem models.CartItems
	result := database.DB.Where("users_id=? AND product_id=?", userId, ProductID).Delete(&removeItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"error": "Product not found",
			})
			return
		} else {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "product remove from the cart",
	})
}

func GetCart(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User id is not exist"})
		return
	}
	var user models.ProfileResponce
	if err := database.DB.Model(&models.Users{}).Where("id=?", userId).First(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "database error"})
		return
	}
	var cartItems []models.CartResponce

	if err := database.DB.Model(&models.CartItems{}).
		Select("cart_items.product_id,cart_items.id,products.product_name, products.image_url, products.description, products.price, products.stock_quantity, products.size, cart_items.quantity,products.discount").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("cart_items.users_id = ?", userId).
		Find(&cartItems).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "database error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"details": "cart item fetched",
		"data": gin.H{
			"user_email": user.Email,
			"cart items": cartItems,
		},
	})

}

func EditCart(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "user id is missing"})
		c.Abort()
		return
	}
	cartItemId := c.Param("id")
	fmt.Println(cartItemId)
	cartItemID, err := strconv.Atoi(cartItemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	type EditCart struct {
		Quantity int `json:"quantity"`
	}
	var Editcart EditCart
	if err := c.BindJSON(&Editcart); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	if Editcart.Quantity <= 0 {
		c.JSON(400, gin.H{"error": "Quantity should be at least 1. The item will not be added to cart"})
		return
	}
	if Editcart.Quantity > 10 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Product quantity exceed maximum limit.Only 10 product quantity availible"})
		return
	}
	var cart models.CartItems

	if err := database.DB.Model(&cart).Where("id=?", cartItemID).Update("quantity", Editcart.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update cart"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusok",
		"details": "cart quantity is updated",
	})

}
