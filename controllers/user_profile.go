package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func GetUserProfile(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "user id is missing"})
		return
	}
	var user models.Users
	var userlist models.ProfileResponce
	result := database.DB.Model(&user).First(&userlist, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"staus": "error",
				"code":  "StatusNotfound(404)",
				"error": "User id is not exist"})
			return
		} else {
			c.JSON(500, gin.H{
				"stauts":  "error",
				"code":    "StatusInternalServerError(500)",
				"message": "Failed to fetch users details",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk(200)",
		"details": "user data fetched succesful",
		"data":    userlist,
	})
}

func EditUserProfile(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "user id is missing"})
		return
	}
	var update_user models.Users
	result := database.DB.First(&update_user, user_id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":   "error",
				"code":     "StatusNotFound(404)",
				"deatails": "User account not found",
			})
			return
		} else {
			c.JSON(500, gin.H{"error": "datatbase error"})
			return
		}
	}
	var input models.ProfileUpdate
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid user request body",
		})
		return
	}

	validate := validator.New()
	validate.RegisterValidation("phone_number", utils.ValidPhoneNum)
	validate.RegisterValidation("alpha_space", utils.ValidateAlphaNumSpace)
	if err := validate.Struct(&input); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    "StatusBadRequest(400)",
			"message": errors,
		})
		return
	}
	updateProfile := make(map[string]interface{})
	if input.FirstName != "" {
		updateProfile["first_name"] = input.FirstName
	}
	if input.LastName != "" {
		updateProfile["last_name"] = input.LastName
	}
	if input.PhoneNumber != "" {
		updateProfile["phone_number"] = input.PhoneNumber
	}
	fmt.Println(updateProfile)

	if err := database.DB.Model(&update_user).Updates(updateProfile).Error; err != nil {
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

func ChangePassword(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "user id is missing",
		})
		c.Abort()
		return
	}
	var changePassword models.ChangePasswordRequest
	if err := c.BindJSON(&changePassword); err != nil {
		c.JSON(400, gin.H{"error": "invlid request"})
		return
	}
	var user models.Users
	result := database.DB.Where("id=?", userId).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not exist",
		})
		return
	}
	validate := validator.New()
	validate.RegisterValidation("password", utils.ValidPassword)
	if err := validate.Struct(&changePassword); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "error",
			"details": "validation error",
			"errors":  errors,
		})
		return
	}

	if err := utils.VerifyPassword(user.Password, changePassword.CurrentPassword); err != nil {
		c.JSON(400, gin.H{"error": "Invalid current password.Please try valid password"})
		return
	}
	HashedPassword, err := utils.HashPassword(changePassword.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Fail to hash password"})
		return
	}
	user.Password = HashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to update new password",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk",
		"details": "user new password updated succesfully",
	})

}

func AddReferralPointsToWallet(c *gin.Context) {
	userId, _ := c.Get("user_id")
	fmt.Printf("%t", userId)
	userID, ok := userId.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	var user models.Users
	if err := database.DB.Model(&models.Users{}).Where("id=?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database errror"})
		return
	}
	ReferalAmount := float64(user.ReferralPoint) / models.ReferalPointsUnit
	
	if user.ReferralPoint < int(models.MinReffaralPoint) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("You need at least %f referral points to convert to cash", models.MinReffaralPoint)})
	}
	

	walletrecord := models.WalleteHistory{
		UserID:         uint(userID),
		Amount:         ReferalAmount,
		Reason:         models.RefferalOffer,
		TransationType: models.CashCredited,
	}
	user.WallteAmount += ReferalAmount
	user.ReferralPoint = 0
	database.DB.Save(&user)

	database.DB.Create(&walletrecord)
	c.JSON(200, gin.H{
		"status":   "success",
		"deatails": "Refferal cash bonus added to the wallete",
		"data": gin.H{
			"wallete_amount": user.WallteAmount,
			"referal_points": user.ReferralPoint,
		},
	})

}

