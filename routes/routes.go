package routes

import (
	middleware "outfiro/middlewares"
    "outfiro/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.Engine) {
	admin := routes.Group("api/v1/admin")
	{
		admin.POST("/singup", controllers.AdminSignup)
		admin.POST("/login", controllers.AdminLogin)
		admin.Use(middleware.AuthMidleware())
		admin.Use(middleware.ProtectedRoutes())
		admin.GET("/users", controllers.GetUsers)
		admin.GET("/users/user/:id", controllers.GetUser)
		admin.PATCH("/users/block/:id", controllers.BlockUser)
		admin.PATCH("/users/unblock/:id", controllers.UnblockUsers)
		admin.GET("/categories", controllers.GetCategories)
		admin.POST("/categories", controllers.AddCategories)
		admin.DELETE("/categories/:id", controllers.DeleteCategory)
		admin.PATCH("/categories/:id", controllers.EditCategory)
		admin.GET("/products", controllers.GetProducts)
		admin.GET("/product/:id", controllers.GetProduct)
		admin.POST("/product", controllers.AddProduct)
		admin.PATCH("/product/:id", controllers.EditProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
	}

}

func UserRoutes(routes *gin.Engine) {
	user := routes.Group("api/v1/user")
	{
		user.POST("/signup", controllers.UserSignup)
		user.POST("/signup/resend-otp", controllers.ResendOtp)
		user.POST("signup/varify-otp", controllers.VerifyOtp)
		user.POST("/login", controllers.UserLogin)
		user.POST("/forgot-password/send-otp", controllers.PasswordSendOtp)
		user.POST("/forgot-password/resend-otp", controllers.ResendOtp)
		user.POST("/forgot-password/verify-otp", controllers.ResetPassOtpVerify)
		user.POST("/forgot-password/reset-password", controllers.ResetPassword)
		user.GET("/google-login", controllers.GoogleLogin)
		user.GET("/google/callback", controllers.GoogleCallback)
		user.Use(middleware.AuthMidleware())
		user.GET("/search", controllers.SearchProduct)
		user.GET("/categories", controllers.GetCategories)
		user.GET("/products", controllers.GetProducts)
		user.GET("/product/:id", controllers.GetProduct)

	}

}
