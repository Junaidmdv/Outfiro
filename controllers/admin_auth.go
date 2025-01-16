package controllers

import (
	"fmt"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"

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
	validate := validator.New()
	if err := validate.RegisterValidation(" ContainsSpecialChars", models.ContainsSpecialChars); err != nil {
		fmt.Println("Erron register custom validation ")
	}
	if err := validate.Struct(&admin); err != nil {
		c.JSON(400, gin.H{
			"statsu": "error",
			"error":  fmt.Sprintf("%v", err),
		})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Admin{}).Where("admin=? or email=?", admin.FirstName, admin.Email).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to check admin existed"})
	}
	if count > 0 {
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
	if err := database.DB.Create(&admin); err != nil {
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
	validate := validator.New()
	if err := validate.Struct(&adminlogin); err != nil {
		c.JSON(400, gin.H{"error": err})
	}
	var admin models.Admin
	result := database.DB.Where("email=?", admin.Email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "email not found"})
		} else {
			c.JSON(500, gin.H{"error": "failed to fetch email details"})
		}
		return
	}
	err := utils.VerifyPassword(admin.Password, adminlogin.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user password"})
		return
	}
	token, err := utils.CreateToken(admin.ID, admin.Email, admin.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create token"})
		return
	}
	c.Header("Authorization", "Bearer"+token)

	c.JSON(200, gin.H{"": "admin login succesful"})

}

func GetUsers(c *gin.Context) {
	var users []models.Users
	if err := database.DB.Select("id", "first_name", "last_name", "email", "address").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch the users"})
	}
	c.JSON(200, gin.H{
		"status":  "error",
		"message": "user fetched succesful",
		"data": gin.H{
			"users": users,
		},
	})
}

func GetUser(c *gin.Context) {
	userIdStr := c.Param("id")
	if userIdStr == " " {
		c.JSON(400, gin.H{"error": "user id is not availible"})
		return
	}
	user_id, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invlid id formate"})
	}
	var user models.Users
	result := database.DB.First(&user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "user id is not availible"})
			return
		} else {
			c.JSON(500, gin.H{"error": "failed fetch id"})
			return
		}
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func BlockUser(c *gin.Context) {
	userIdStr := c.Param("id")
	if userIdStr == " " {
		c.JSON(400, gin.H{"error": "user id is not availible"})
		return
	}
	user_id, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invlid id formate"})
	}
	var user models.Users
	result := database.DB.First(&user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "user id is not availible"})
			return
		} else {
			c.JSON(500, gin.H{"error": "failed fetch id"})
			return
		}
	}
	user.IsBlocked = true
	if err := database.DB.Save(&user); err != nil {
		c.JSON(500, gin.H{"error": "failed to block the user"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "user account blocked",
	})

}

func UnblockUsers(c *gin.Context) {
	userIdStr := c.Param("id")
	if userIdStr == " " {
		c.JSON(400, gin.H{"error": "user id is not availible"})
		return
	}
	user_id, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invlid id formate"})
	}
	var user models.Users
	result := database.DB.First(&user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "user id is not availible"})
			return
		} else {
			c.JSON(500, gin.H{"error": "failed fetch id"})
			return
		}
	}
	user.IsBlocked = false
	if err := database.DB.Save(&user); err != nil {
		c.JSON(500, gin.H{"error": "failed to block the user"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "user account unblocked",
	})

}
