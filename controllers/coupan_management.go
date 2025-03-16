package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CreateCoupan(c *gin.Context) {
	var couponRequest models.Coupon
	if err := c.BindJSON(&couponRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body %v", err)})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&couponRequest); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"errors": errors,
		})
		return
	}
	var exist bool
	if err := database.DB.
		Model(&models.Coupon{}).
		Select("COUNT(*)>0").
		Where("coupon_code=?", couponRequest.CouponCode).
		Scan(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if exist {
		c.JSON(http.StatusConflict, gin.H{"error": "coupan already exist"})
		return
	}

	if err := database.DB.Create(&couponRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"details": "new coupan created",
		"data":    couponRequest,
	})

}

func DeleteCoupan(c *gin.Context) {
	couponId := c.Param("id")
	couponID, err := strconv.Atoi(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupan id"})
		return
	}

	if err := database.DB.Delete(models.Coupon{}, couponID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "error",
		"details": "Coupon deleted",
	})
}

func ListCoupan(c *gin.Context) {
	var Coupons []models.Coupon

	if err := database.DB.Find(&Coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch coupon",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"details": "coupan data fetched",
		"data":    Coupons,
	})
}

func EditCoupan(c *gin.Context) {
	couponId := c.Param("id")
	couponID,_:=strconv.Atoi(couponId)
	var coupon models.Coupon
	result:=database.DB.Model(&coupon).First(&coupon, couponID)
	if result.Error !=nil{
		if errors.Is(result.Error,gorm.ErrRecordNotFound){
			 c.JSON(http.StatusNotFound,gin.H{"error":"Coupon not found"})
			 return
		}else{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"database error"})
			return
		}
	}
	var updateCoupan models.CouponUpdate
	if err := c.BindJSON(&updateCoupan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	fmt.Println(updateCoupan.MinPurchase)
	validate := validator.New()
	if err := validate.Struct(&updateCoupan); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"errors": errors,
		})
		return
	}
	update := make(map[string]interface{})
	if updateCoupan.DiscountPercent != 0 {
		update["discount_percent"] = updateCoupan.DiscountPercent
	}
	if updateCoupan.MinPurchase != 0 {
		update["min_purchase"] = updateCoupan.MinPurchase
	}
	

	if err := database.DB.Model(&coupon).Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update the coupon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"details": "coupon updated",
	})

}
