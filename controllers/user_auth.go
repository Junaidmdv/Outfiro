package controllers

import (
	"errors"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrEmailExist     = errors.New("email already existed")
	ErrPasswordLength = errors.New("password should be 8 charectar long")
	ErrInternalServer = errors.New("internal server error")
)

func UserSignup(c *gin.Context) {
	var usersignup models.SignuPlayload
	if err := c.BindJSON(&usersignup); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalide code",
		})
		return
	}
	if result := utils.IsEmailValid(usersignup.Email); !result {
		c.JSON(400, gin.H{
			"error": "Invalid email formate",
		})
		return
	}

	if len(usersignup.Password) < 8 {
		c.JSON(400, gin.H{
			"error": ErrPasswordLength,
		})
		return
	}
	if usersignup.ConfirmPassword != usersignup.Password {
		c.JSON(400, gin.H{
			"error": "password and confirm password do not match",
		})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Users{}).Where("email = ?", usersignup.Email).Count(&count); err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	if count > 0 {
		c.JSON(400, ErrEmailExist)
		return
	}
	HashedPassword, err := utils.HashPassword(usersignup.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"errors": err,
		})
		return
	}
	usersignup.Password = HashedPassword
	var otp utils.OTP
	err = otp.GenerateOTP(6)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	email := strings.Split(usersignup.Email, ",")
	err = utils.SendEmail(email, otp.Value)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	otp.Email = email[0]
	if err := database.DB.Create(&usersignup); err != nil {
		c.JSON(500, gin.H{"error": "failed to add the otp details"})
	}
	if err := database.DB.Create(&otp); err != nil {
		c.JSON(500, gin.H{"error": "failed to add the otp details"})
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "user signup successful",
	})
}

func VerifyOtp(c *gin.Context) {
	var verifyotp models.OtpVerifier
	if err := c.BindJSON(&verifyotp); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid code formate ",
		})
		return
	}
	var otp utils.OTP
	if err := otp.Verifyotp(verifyotp.Value); err != nil {
		c.JSON(400, gin.H{"errors": err})
	}
	var user models.Users
	var SingnupData models.SignuPlayload
	if err := database.DB.Find("email", otp.Email).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to sigup details",
		})
	}

	user.FirstName = SingnupData.FirstName
	user.LastName = SingnupData.LastName
	user.Email = SingnupData.Email
	user.Password = SingnupData.Password

	if err := database.DB.Create(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed add user data"})
	}
	c.JSON(200, gin.H{
		"status": "success",
		"user":   user,
	})

}

func UserLogin(c *gin.Context) {

	var userlogin models.LoginPlayload

	if err := c.BindJSON(&userlogin); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid data",
		})
		c.Abort()
		return
	}
	var userData models.Users

	result := database.DB.Where("name=? and ", userlogin.Email).First(&userData)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "fail to found the user deatail"})
		} else {
			c.JSON(500, gin.H{"error": "internal server error"})
		}
		return
	}

	if err := utils.VerifyPassword(userlogin.Password, userData.Password); err != nil {
		c.JSON(401, gin.H{
			"error": "invalid user password",
		})
		return
	}

	token, err := utils.CreateToken(userData.ID, userData.Email, userData.Role)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create token",
		})
	}
	c.Header("Authorization", "Bearer"+token)

	c.JSON(200, gin.H{
		"status":  "succesful",
		"message": "user login succesful",
		"token":   token,
	})

}
