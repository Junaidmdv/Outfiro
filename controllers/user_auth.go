package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	ErrEmailExist     = errors.New("email already existed")
	ErrPasswordLength = errors.New("password should be 8 charectar long")
	ErrInternalServer = errors.New("internal server error")
	OtpExpirePeriod   = 5
)

func UserSignup(c *gin.Context) {
	var usersignup models.SignuPlayload
	var user models.Users
	if err := c.BindJSON(&usersignup); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    400,
			"details": "Field validation error",
		})
		return
	}
	validate := validator.New()
	validate.RegisterValidation("phone_number", utils.ValidPhoneNum)
	validate.RegisterValidation("password", utils.ValidPassword)
	if err := validate.Struct(&usersignup); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    "StatusBadRequest(400)",
			"details": "validation error occur",
			"errors":  errors,
		})
		return
	}

	if usersignup.ConfirmPassword != usersignup.Password {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    400,
			"details": "password and confirm password do not match",
		})
		return
	}
	var userExist bool
	err := database.DB.Model(&user).
		Select("count(*) > 0").
		Where("email=?", usersignup.Email).
		Scan(&userExist).Error
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    500,
			"details": "Unexpected error occur on processing the request",
		})
		return
	}
	fmt.Println(userExist)
	if userExist {
		c.JSON(409, gin.H{
			"status":  "error",
			"code":    "409",
			"details": "This email is already registered. Please log in or use a different email to sign up",
		})
		return
	}

	var NumExist bool
	err = database.DB.Model(&user).
		Select("count(*) > 0").
		Where("phone_number=?", usersignup.PhoneNumber).
		Scan(&NumExist).Error
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    500,
			"details": "Unexpected error occur on processing the request",
		})
		return
	}
	fmt.Println(NumExist)
	if NumExist {
		c.JSON(409, gin.H{
			"status":  "error",
			"code":    "409",
			"details": "This phone number is already registered. Please log in or use a different number to sign up",
		})
		return
	}
	var AdminExist bool
	var admin models.Admin
	if err := database.DB.Model(&admin).
		Select("count(*) > 0").
		Where("email=?", usersignup.Email).
		Scan(&AdminExist).Error; err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"code":     500,
			"deatails": "Unexpected error occur on processing the request",
		})
		return
	}
	fmt.Println(userExist)
	if AdminExist {
		c.JSON(409, gin.H{
			"status":  "error",
			"code":    "409",
			"details": "The entered email belongs to an admin account and cannot be used for user registration",
		})
		return
	}
	//validation for referal code
	if usersignup.ReferedCode != "" {
		var exist bool
		if err := database.DB.Model(models.Users{}).Select("COUNT(*)>0").
			Where("referral_code=?", usersignup.ReferedCode).
			Scan(&exist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		if !exist {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid referal code"})
			return
		}

	}

	HashedPassword, err := utils.HashPassword(usersignup.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    500,
			"details": err.Error(),
		})
		return
	}
	HashConPassword, err := utils.HashPassword(usersignup.ConfirmPassword)
	if err != nil {
		c.JSON(500, gin.H{
			"errors": err.Error(),
		})
		return
	}

	var otp utils.OTP
	otp.Email = usersignup.Email
	if err := otp.GenerateOTP(6); err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"code":     500,
			"deatails": "Failed to generate otp",
		})
		return
	}
	email := strings.Split(otp.Email, ",")
	if err := utils.SendEmail(email, otp.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to share otp",
		})
		return
	}
	var otpRecord models.OtpRecord
	otpRecord.Value = otp.Value
	otpRecord.Email = otp.Email
	otpRecord.ExpiryTime = time.Now().Add(time.Minute * time.Duration(OtpExpirePeriod))

	if err := database.DB.Create(&otpRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to add the otp data %v",err.Error()),
		})
		return
	}

	usersignup.Password = HashedPassword
	usersignup.ConfirmPassword = HashConPassword
	if err := database.DB.Create(&usersignup).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to add the user details"})
		return
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"code":     200,
		"deatails": "user signup successful",
	})
}

func ResendOtp(c *gin.Context) {
	var req utils.OTP
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid inpute formate",
		})
		return
	}
	var ResendOtp models.OtpRecord
	result := database.DB.Where("email=?", req.Email).Find(&ResendOtp)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "email is not matching"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database failure to check the user email"})
			return
		}
	}
	fmt.Println(ResendOtp)
	if time.Now().Before(ResendOtp.ExpiryTime) {
		c.JSON(400, gin.H{
			"error": "otp already sent",
		})
		return
	}
	if err := req.GenerateOTP(6); err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"code":     500,
			"deatails": "Failed to generate otp",
		})
		return
	}
	email := strings.Split(req.Email, ",")
	if err := utils.SendEmail(email, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to share otp",
		})
		return
	}
	ResendOtp.ExpiryTime = time.Now().Add(time.Minute * time.Duration(OtpExpirePeriod))
	ResendOtp.Value = req.Value
	ResendOtp.Purpose = "Resend Otp"
	if err := database.DB.Save(&ResendOtp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in the database"})
		return
	}
	c.JSON(200, gin.H{
		"status":   "success",
		"deatails": "Otp sends succesful",
	})
}

func VerifyOtp(c *gin.Context) {
	var verifyotp utils.OTP
	if err := c.BindJSON(&verifyotp); err != nil {
		c.JSON(400, gin.H{
			"status":   "error",
			"code":     400,
			"deatails": "Invalid code formate ",
		})
		return
	}
	var otp models.OtpRecord
	result := database.DB.Where("email= ?", verifyotp.Email).Order("created_at desc").First(&otp)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":   "error",
				"code":     404,
				"deatails": "otp data is not availible",
			})
		} else {
			c.JSON(500, gin.H{
				"status":   "error",
				"code":     "500",
				"deatails": "failed to fetch opt data"})
		}
		return
	}
	if err := verifyotp.Verifyotp(&otp); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":   "error",
			"code":     "StatusUnauthorized(401)",
			"deatails": err.Error(),
		})
		return
	}
	var signupData models.SignuPlayload
	if err := database.DB.Where("email", verifyotp.Email).First(&signupData).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to take the signup details"})
		return
	}
	user := models.Users{
		FirstName:    signupData.FirstName,
		LastName:     signupData.LastName,
		PhoneNumber:  signupData.PhoneNumber,
		Email:        signupData.Email,
		Password:     signupData.Password,
		ReferredCode: signupData.ReferedCode,
	}
	user.ReferralCode = utils.GenerateReferralCode()

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//adding referral point to the referred user
	if signupData.ReferedCode != "" {
		result = database.DB.Model(&models.Users{}).
			Where("referral_code=?", signupData.ReferedCode).
			Update("referral_point", gorm.Expr("referral_point + ?", models.ReferralPoints))

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update the referal points"})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to update the referal points"})
			return
		}
	}

	if err := database.DB.Unscoped().Where("email", verifyotp.Email).Delete(&signupData).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "error delete the user signup data",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk(200)",
		"details": "new user account created",
	})

}

func UserLogin(c *gin.Context) {
	var userlogin models.LoginPlayload
	if err := c.BindJSON(&userlogin); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		c.Abort()
		return
	}
	var userData models.Users
	result := database.DB.Where("email=?", userlogin.Email).First(&userData)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"code":    "StatusNotFound(404)",
				"details": "User account not found",
			})
		} else {
			c.JSON(500, gin.H{"error": "failed to check the user details"})
		}
		return
	}
	if userData.IsBlocked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":   "error",
			"code":     401,
			"deatails": "user account is blocked",
		})
		return
	}
	if err := utils.VerifyPassword(userData.Password, userlogin.Password); err != nil {
		c.JSON(401, gin.H{
			"error": "invalid user password",
		})
		return
	}
	fmt.Println(userData)
	token, err := utils.GenerateToken(userData.ID, userData.Email, userData.Role)
	if err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"deatails": "Failed to create token",
			"code":     "ErrInternalServer(500)",
		})
	}
	c.Header("Authorization", token)

	c.JSON(200, gin.H{
		"status":   "succesful",
		"deatails": "user login succesful",
		"token":    token,
	})
}

func GoogleLogin(c *gin.Context) {
	fmt.Println("Hellow google")
	GoogleAuthConfig := utils.GoogleConfig()
	url := GoogleAuthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)

}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "state" {
		c.JSON(400, gin.H{
			"status":   "error",
			"deatails": "state don't Match!!",
		})
		return
	}

	code := c.Query("code")
	googleConfig := utils.GoogleConfig()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(400, gin.H{
			"status":   "error",
			"deatails": "Failed to get the user details",
		})
		return
	}
	responce, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"deatails": "failed pass the user details",
		})
		return
	}
	userData, err := io.ReadAll(responce.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"deatails": "json parsing failed",
		})
		return
	}
	var newUser models.GoogleUser

	if err := json.Unmarshal(userData, &newUser); err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"deatails": "failed to fetch user details",
		})
		return
	}
	user := models.Users{
		FirstName: newUser.Name,
		LastName:  newUser.FamilyName,
		Email:     newUser.Email,
		Password:  " ",
		Role:      "user",
		LoginType: "google",
	}
	var exist bool
	if database.DB.Model(user).Where("email", newUser.Email).Scan(&exist); exist {
		fmt.Println("user email already exist")

	} else {
		err = database.DB.Model(&user).Create(&user).Error
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed add the details in the database",
			})
			return
		}
	}
	var CreatedUser models.Users
	if err := database.DB.Where("email=?", newUser.Email).First(&CreatedUser); err != nil {
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	Authtoken, err := utils.GenerateToken(CreatedUser.ID, CreatedUser.Email, CreatedUser.Role)
	if err != nil {
		c.JSON(500, gin.H{
			"deatails": "failed to create authentication token",
		})
		return
	}
	c.Header("Authorization", Authtoken)
	c.JSON(200, gin.H{
		"status":   "success",
		"deatails": "user details fetched succesful",
		"token":    Authtoken,
	})

}
