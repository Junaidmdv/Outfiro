package controllers

import (
	"errors"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func AddAddress(c *gin.Context) {
	UserId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing"})
		return
	}
	var user models.Users
	result := database.DB.Model(&user).Where("id=?", UserId).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"code":    "StatusNotFound(404)",
				"details": "User account not found",
			})
			return
		} else {
			c.JSON(500, gin.H{"error": "datatbase error"})
			return
		}
	}
	var address models.Address
	if err := c.BindJSON(&address); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	address.UsersID = user.ID
	validate := validator.New()
	validate.RegisterValidation("pincode", utils.ValidPincode)
	validate.RegisterValidation("alpha_space", utils.ValidateAlphaNumSpace)
	if err := validate.Struct(&address); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "erros",
			"code":    "StatusBadRequest",
			"details": "validation error",
			"errors":  errors,
		})
		return
	}
	if err := database.DB.Create(&address).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed add user address"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"satus":   "success",
		"code":    "StatusCreated(201)",
		"details": "New address added succesful",
	})

}

func EditAddress(c *gin.Context) {
	UserId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing"})
		return
	}
	AddIdstr := c.Param("id")
	AddId, _ := strconv.Atoi(AddIdstr)

	var userAdd models.Address
	result := database.DB.Model(&userAdd).Where("id=? and users_id=?", AddId, UserId).First(&userAdd)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"code":    "StatusNotFound(404)",
				"details": "User account not found",
			})
			return
		} else {
			c.JSON(500, gin.H{"error": "datatbase error"})
			return
		}
	}
	var updateAdd models.AddressResponce
	if err := c.BindJSON(&updateAdd); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	validate := validator.New()
	validate.RegisterValidation("pincode", utils.ValidPincode)
	validate.RegisterValidation("alpha_space", utils.ValidateAlphaNumSpace)
	if err := validate.Struct(&updateAdd); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "erros",
			"code":    "StatusBadRequest",
			"details": "validation error",
			"errors":  errors,
		})
		return
	}
	UpdateAddress := make(map[string]interface{})
	if updateAdd.Address_line_1 != "" {
		UpdateAddress["address_line_1"] = updateAdd.Address_line_1
	}
	if updateAdd.Address_line_2 != "" {
		UpdateAddress["address_line_2"] = updateAdd.Address_line_2
	}
	if updateAdd.City != "" {
		UpdateAddress["city"] = updateAdd.City
	}
	if updateAdd.Contact_number != "" {
		UpdateAddress["contact_number"] = updateAdd.Contact_number
	}
	if updateAdd.Country != " " {
		UpdateAddress["country"] = updateAdd.Country
	}
	if updateAdd.Landmark != "" {
		UpdateAddress["pincode"] = updateAdd.Pincode
	}
	if updateAdd.State != "" {
		UpdateAddress["state"] = updateAdd.State
	}

	if err := database.DB.Model(&userAdd).Updates(updateAdd).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    "StatusInternalServerError(500)",
			"details": "failed update the user details",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "StatusOk",
		"message": "User details are updated",
	})
}

func DeleteAddress(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing"})
		return
	}
	AddressIdStr := c.Param("id")
	AddressId, err := strconv.Atoi(AddressIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid address id"})
	}
	var userAdd models.Address
	result := database.DB.Where("id=? and users_id=?", AddressId, userId).First(&userAdd)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"code":    "StatusNotFound(404)",
			"details": "User address not found found",
		})
		return
	}
	if err := database.DB.Delete(&userAdd).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to delete address"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk",
		"details": "user acccount deleted succesful",
	})

}

func GetAddress(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing"})
		return
	}
	var Add models.Address
	var Address []models.AddressResponce
	result := database.DB.Model(&Add).Where("users_id=?", userId).Find(&Address)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    "StatusNotFound",
				"details": "User id is not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}
	var user models.Users
	if err := database.DB.Where("id=?", userId).First(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "datatbaser error"})
		return
	}
	c.JSON(200, gin.H{
		"staus": "success",
		"code":  "statusOk",
		"data": gin.H{
			"user":    user.Email,
			"address": Address,
		},
	})

}
