package controllers

import (
	"fmt"
	"outfiro/database"
	"outfiro/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context) {
	var userlist []models.GetUser
	var users models.Users
	if err := database.DB.Model(&users).Find(&userlist).Error; err != nil {
		c.JSON(500, gin.H{
			"stauts":  "error",
			"code":    "StatusInternalServerError(500)",
			"message": "Failed to fetch user details",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "error",
		"code":    200,
		"message": "user fetched succesful",
		"data": gin.H{
			"users": userlist,
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
	var userlist models.GetUser
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
		"status": "success",
		"data": gin.H{
			"user": userlist,
		},
	})
}

func BlockUser(c *gin.Context) {
	userIdStr := c.Param("id")
	if userIdStr == " " {
		c.JSON(400, gin.H{"error": "user id is not availible"})
		return
	}
	fmt.Println("userIdStr")
	user_id, err := strconv.Atoi(userIdStr)
	fmt.Println(user_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invlid id formate"})
		return
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
	if !user.IsBlocked {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    400,
			"message": "User already in the blocked list",
		})
		return
	}
	user.IsBlocked = false

	if err := database.DB.Model(&user).Update("is_blocked", false).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to block the user"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "user account unblocked",
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
	if user.IsBlocked {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    400,
			"message": "User already in the unblocked list",
		})
		return
	}

	if err := database.DB.Model(&user).Update("is_blocked", false).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to block the user"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"code":    200,
		"message": "user account unblocked",
	})

}

// func AddUser(c *gin.Context) {
// 	var user models.Users
// 	if err := c.Bind(&user); err != nil {
// 		c.JSON(400, gin.H{
// 			"status":  "error",
// 			"code":    400,
// 			"message": "Invalid input formate",
// 		})
// 		return
// 	}
// 	validate := validator.New()
// 	validate.RegisterValidation("containsSpecialChars", models.ContainsSpecialChars)

// 	if err := validate.Struct(&user); err != nil {
// 		errors := utils.UserFormateError(err.(validator.ValidationErrors))
// 		c.JSON(400, gin.H{
// 			"status": "errors",
// 			"code":   400,
// 			"errors": errors,
// 		})
// 		return
// 	}

// 	HashedPassword, err := utils.HashPassword(user.Password)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"status":"error",
// 			"code":500,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	user.Password=HashedPassword
// 	var exist bool

// 	if err := database.DB.Model(user).
// 		Select("count(*) > 0").
// 		Where("email=?", user.Email).
// 		Scan(&exist).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"status":  "error",
// 			"code":    500,
// 			"message": "Unexpected error occur on processing the request",
// 		})
// 		return

// 	}
// 	if exist {
// 		c.JSON(400, gin.H{
// 			"status":  "error",
// 			"code":    400,
// 			"message": "User already exist in the database",
// 		})
// 		return
// 	}

// 	if err := database.DB.Create(&user).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"status":  "error",
// 			"code":    400,
// 			"message": "Failed to add user in database",
// 		})
// 		return
// 	}
// 	c.JSON(201, gin.H{
// 		"status":  "success",
// 		"code":    "StatusCreated",
// 		"message": "new user is created",
// 	})

// }

// // func UpdateUserDetails(c *gin.Context) {
// // 	var update_user models.Users
// // 	userIdstr := c.Param("id")
// // 	if userIdstr == " " {
// // 		c.JSON(400, gin.H{
// // 			"status":  "error",
// // 			"code":    "StatusBadRequest(400)",
// // 			"message": "query parameters are required",
// // 		})
// // 	}

// // 	userId, err := strconv.Atoi(userIdstr)
// // 	if err != nil {
// // 		c.JSON(400, gin.H{
// // 			"status":  "error",
// // 			"code":    "StatusBadRequest(400)",
// // 			"message": "Invalid query parameters",
// // 		})
// // 	}

// // 	c.BindJSON(&update_user)
// // 	validate := validator.New()
// // 	if err := validate.Struct(&update_user); err != nil {
// // 		errors := utils.UserFormateError(err.(validator.ValidationErrors))
// // 		c.JSON(400, gin.H{
// // 			"status":  "error",
// 			"code":    "StatusBadRequest(400)",
// 			"message": errors,
// 		})
// 		return
// 	}

// 	result := database.DB.Model(&models.Users{}).Where("id=?", userId).Updates(update_user)

// 	if result.Error != nil {
// 		c.JSON(500, gin.H{
// 			"status":  "error",
// 			"code":    "StatusInternalServerError(500)",
// 			"message": "Failed to update user",
// 		})
// 		return
// 	}
// 	if result.RowsAffected == 0 {
// 		c.JSON(404, gin.H{
// 			"status":  "error",
// 			"code":    "StatusNotfound(404)",
// 			"message": "user not found",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"status":  "success",
// 		"code":    "StatusOk",
// 		"message": "User details are updated",
// 	})
// }

// func DeleteUser(c *gin.Context) {
// 	var deleteUser models.Users
// 	userIdstr := c.Param("id")
// 	if userIdstr == " " {
// 		c.JSON(400, gin.H{
// 			"status":  "error",
// 			"code":    "BadRequest(400)",
// 			"message": "User id is missing",
// 		})
// 	}
// 	userId, err := strconv.Atoi(userIdstr)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"status":  "error",
// 			"code":    "BadRequest(400)",
// 			"message": "Invalid user id",
// 		})
// 	}
// 	var exist bool
// 	if err := database.DB.Model(&deleteUser).
// 		Select("count(*)>0").Where("id=?", userId).
// 		Scan(exist).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"status":  "error",
// 			"code":    "InternalServer(500)",
// 			"message": "database error occured",
// 		})
// 	}
// 	if !exist {
// 		c.JSON(400, gin.H{
// 			"status":  "error",
// 			"code":    "BadRequest",
// 			"message": "User id is missing",
// 		})
// 	}

// 	if err := database.DB.Delete(&deleteUser, userId).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"status":  "error",
// 			"code":    "InternalServer(500)",
// 			"message": "database error occured",
// 		})
// 	}

// 	c.JSON(200, gin.H{
// 		"status":  "success",
// 		"code":    "StatusOk",
// 		"message": "user acount deleted succesful",
// 	})
// }
