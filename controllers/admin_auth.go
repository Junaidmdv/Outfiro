package controllers

import (
	"fmt"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func AdminSignup(c *gin.Context) {
	var admin models.Admin
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(400, gin.H{"error": "invalid input formate"})
		return
	}
	fmt.Println(admin)
	validate := validator.New()
	validate.RegisterValidation("password", utils.ValidPassword)
	if err := validate.Struct(&admin); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status": "error",
			"error":  errors,
		})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Admin{}).Where("email=?", admin.Email).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to check admin existed"})
		return
	}
	if count != 0 {
		c.JSON(409, gin.H{"error": "admin already exist"})
		return
	}
	HashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"errors": fmt.Sprintf("Error in the hashing password %v", err),
		})
		return
	}
	admin.Password = HashedPassword
	if err := database.DB.Create(&admin).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to signup"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "admin created succesfully",
	})
}

func AdminLogin(c *gin.Context) {
	var adminlogin models.AdminLogin
	if err := c.BindJSON(&adminlogin); err != nil {
		c.JSON(400, gin.H{"error": "invalid input formate"})
	}
	var admin models.Admin
	result := database.DB.Where("email=?", adminlogin.Email).First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "email not found"})
		} else {
			c.JSON(500, gin.H{"error": "failed to fetch email details"})
		}
		return
	}
	fmt.Println(admin.Password, adminlogin.Password)
	if err := utils.VerifyPassword(admin.Password, adminlogin.Password); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user password"})
		return
	}

	SignedToken, err := utils.GenerateToken(admin.ID, admin.Email, admin.Role)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to create token",
			"code":    "ErrInternalServer(500)",
		})
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User login succesful",
		"code":    "statusOk(200)",
		"token":   SignedToken,
	})

}
