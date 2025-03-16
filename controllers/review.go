package controllers

import (
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Addreview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userId := userID.(int)

	ProductId := c.Param("id")
	ProductID, err := strconv.Atoi(ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var reviewReq models.ReviewRequest

	if err := c.BindJSON(&reviewReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	validater := validator.New()
	if err := validater.Struct(reviewReq); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}
	var exist bool
	if err := database.DB.Model(&models.Review{}).
		Select("COUNT(*) > 0").
		Where("product_id=? AND user_id=?", ProductID, userID).
		Scan(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{"error": "Product review already added"})
		return

	}
	review := models.Review{
		UserID:    uint(userId),
		ProductID: uint(ProductID),
		Rating:    reviewReq.Rating,
		Comments:  reviewReq.Comments,
	}
	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"details": "new data is created",
	})
}

func UpdateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ProductId := c.Param("id")
	ProductID, err := strconv.Atoi(ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	var exist bool
	if err := database.DB.Model(&models.Review{}).
		Select("COUNT(*) > 0").
		Where("product_id=? AND user_id=?", ProductID, userID).
		Scan(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exist {
		c.JSON(http.StatusConflict, gin.H{"error": "Product not exist"})
		return
	}

	var reviewReq models.ReviewRequest
	if err := c.BindJSON(&reviewReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	updateReview := make(map[string]any)
	if reviewReq.Rating != 0 {
		updateReview["rating"] = reviewReq.Rating
	}
	if reviewReq.Comments != "" {
		updateReview["comments"] = reviewReq.Comments
	}
	fmt.Println(updateReview)
	if err := database.DB.Model(&models.Review{}).Where("user_id=? AND product_id=?", userID, ProductID).Updates(&updateReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":  "success",
		"details": "Review Updated",
	})

}

func DeleteReview(c *gin.Context) {
	ProductID, err := strconv.Atoi(c.Param("id"))
	fmt.Println(ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Product id"})
		return
	}

	var review models.Review
	if err := database.DB.Where("product_id=?", ProductID).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product review not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if review.UserID != uint(userID.(int)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to delete this review"})
		return
	}

	if err := database.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Review deleted"})
}

func GetProductReview(c *gin.Context) {
	productId := c.Param("id")
	productID, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productId"})
	}

	var exist bool
	if err := database.DB.Model(&models.Products{}).
		Select("COUNT(*) > 0").
		Where("id=?", productID).
		Scan(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "products are not found in the database"})
		return
	}

	var review []models.Review
	if err := database.DB.Where("product_id=?", productID).Preload("User").Find(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if len(review) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"details": "product review is empty",
		})
	}
	var avgRating float64
	var RatingCount uint
	var Rating float64
	var revResponce []models.ReviewResponce
	for _, items := range review {

		revResponce = append(revResponce, models.ReviewResponce{
			ID:         items.ID,
			UserName:   fmt.Sprintf("%s %s", items.User.FirstName, items.User.LastName),
			Rating:     items.Rating,
			Comments:   items.Comments,
			Created_at: items.CreatedAt.Format("2006-01-02 15:04:05"),
		})

		if items.Rating != 0 {
			Rating += items.Rating
			RatingCount++
		}
	}
	avgRating = Rating / float64(RatingCount)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"product_id":     productID,
			"average_rating": avgRating,
			"rate_count":     RatingCount,
			"product_review": revResponce,
		},
	})

}
