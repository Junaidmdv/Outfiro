package controllers
import ("github.com/gin-gonic/gin"
       "github.com/go-playground/validator/v10"
	   "net/http"
	    "outfiro/utils"
		"outfiro/models"
		 "outfiro/database"
		 "gorm.io/gorm"
		 "time"
		 "strings"
		 "fmt"
)

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
	if err := database.DB.Exec("UPDATE users SET password = ? Where email=?", HashPassword, otp.Email).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    "StatusInternalServerError(500)",
			"details": "Failed update data",
		})
		return
	}

	if err := database.DB.Exec("UPDATE otp_records SET varify = ? Where email=?", false, otp.Email).Error; err != nil {
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
