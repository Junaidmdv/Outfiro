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
	OtpExpirePeriod   = 2
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
	if result := utils.IsEmailValid(usersignup.Email); !result {
		c.JSON(400, gin.H{
			"error": "Invalid email formate",
		})
		return
	}

	if len(usersignup.Password) < 8 {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    400,
			"details": "Password must be atleast 8 charectar long",
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
			"details": "User already exist",
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

func SendOtp(c *gin.Context) {
	var otpRequest utils.OTP
	if err := c.BindJSON(&otpRequest); err != nil {
		c.JSON(400, gin.H{
			"status":   "error",
			"code":     400,
			"deatails": "invalid code",
		})
		return
	}
	var UserDetails models.SignuPlayload
	result := database.DB.Model(&UserDetails).Where("email=?", otpRequest.Email)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"Status":   "error",
				"code":     "StatusNotFound",
				"deatails": "Invalid user email",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "error",
				"code":     "StatusInternalServerError(500)",
				"deatails": "Database Failure",
			})
			return
		}
	}
	if err := otpRequest.GenerateOTP(6); err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"code":     500,
			"deatails": "Failed to generate otp",
		})
		return
	}
	email := strings.Split(otpRequest.Email, ",")
	if err := utils.SendEmail(email, otpRequest.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to share otp",
		})
		return
	}
	var otp models.OtpRecord
	otp.Value = otpRequest.Value
	otp.Email = otpRequest.Email
	otp.ExpiryTime = time.Now().Add(time.Minute * time.Duration(OtpExpirePeriod))

	if err := database.DB.Create(&otp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add the otp data",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"code":     200,
		"deatails": "otp sent succesfully",
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
	// if err := database.DB.Update("verify=?", true).Error; err != nil {
	// 	c.JSON(500, gin.H{
	// 		"error": "Failed verfiy the otp.",
	// 	})
	// 	return
	// }
	var signupData models.SignuPlayload
	if err := database.DB.Where("email", verifyotp.Email).First(&signupData).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to take the signup details"})
		return
	}
	user := models.Users{
		FirstName: signupData.FirstName,
		LastName:  signupData.LastName,
		Email:     signupData.Email,
		Password:  signupData.Password,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Unscoped().Where("email", verifyotp.Email).Delete(&signupData).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "error delete the user signup data",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"code":"statusOk(200)",
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
				"status": "error",
				"code":"StatusNotFound(404)",
				"details":"User account not found",
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

	token, err := utils.GenerateToken(userData.Email, userData.Role)
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
		fmt.Println("user exist")

	} else {
		err = database.DB.Model(&user).Save(&user).Error
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed add the details in the database",
			})
			return
		}
	}

	Authtoken, err := utils.GenerateToken(newUser.Email, "user")
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

func PasswordSendOtp(c *gin.Context) {
	var input utils.OTP
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	validater := validator.New()
	if err := validater.Struct(&input); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, gin.H{
			"status":   "error",
			"code":     "StatusBadRequest(400)",
			"deatails": errors,
		})
		return
	}
	var user models.Users
	result := database.DB.Model(&user).Where("email=?", input.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":   "error",
				"code":     "StatusNotFound(404)",
				"deatails": "user is not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	if err := input.GenerateOTP(6); err != nil {
		c.JSON(500, gin.H{
			"status":   "error",
			"code":     500,
			"deatails": "Failed to generate otp",
		})
		return
	}
	email := strings.Split(input.Email, ",")
	fmt.Println(email, input.Value)
	fmt.Println(email)
	if err := utils.SendEmail(email, input.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to share otp",
		})
		return
	}
	var otp models.OtpRecord
	otp.Value = input.Value
	otp.Email = input.Email
	otp.ExpiryTime = time.Now().Add(time.Minute * 1)
	otp.Purpose = "forgotPassword"

	if err := database.DB.Create(&otp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add the otp data",
		})
		return
	}
	c.JSON(201, gin.H{
		"status":  "success",
		"details": "otp send succesful",
		"code":    "StatusCreated(201)",
	})

}

func ResetPassOtpVerify(c *gin.Context) {
	var request utils.OTP
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"status":   "error",
			"deatails": "invalid input formate",
			"code":     "StatusBadRequest(400)",
		})
		return
	}
	var otp models.OtpRecord
	result := database.DB.Where("email=?", request.Email).First(&otp)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":   "error",
				"deatails": "email is not found",
				"code":     "StatusNotFound(404)",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"status":   "error",
				"deatails": "record not found",
			})
			return
		}
	}
	if err := request.Verifyotp(&otp); err != nil {
		c.JSON(401, gin.H{
			"status":   "error",
			"code":     "StatusUnauthorized",
			"deatails": err.Error(),
		})
		return
	}
	if err := database.DB.Model(&otp).Update("varify", true).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed verfiy the user",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "otp verification success",
	})
}

func ResetPassword(c *gin.Context) {
	var newPassword models.ResetPassword
	if err := c.BindJSON(&newPassword); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}
	fmt.Println("hello")
	var otp models.OtpRecord
	result := database.DB.Where("email=?", newPassword.Email).First(&otp)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{
				"error": "Failed the fetch user email",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"error": "database failure",
			})
			return
		}
	}
	if !otp.Varify {
		c.JSON(401, gin.H{
			"status":  "error",
			"code":    "StatusUnauthorized(401)",
			"details": "otp verfication failed.Please try the otp verification",
		})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&newPassword); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":   "error",
			"code":     "StatusBadRequest(400)",
			"deatails": errors,
		})
		return
	}
	HashPassword, err := utils.HashPassword(newPassword.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "error in password hashing"})
	}
	if err := database.DB.Exec("UPDATE users SET password = ? Where email=?", HashPassword,otp.Email).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    "StatusInternalServerError(500)",
			"details": "Failed update data",
		})
		return
	}

	if err := database.DB.Exec("UPDATE otp_records SET varify = ? Where email=?",false,otp.Email).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    "StatusInternalServerError(500)",
			"details": "Failed update data",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    "statusOk(200)",
		"details": "reset password success",
	})
}
