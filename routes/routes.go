package routes

import (
	"outfiro/controllers"
	middleware "outfiro/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.Engine) {
	admin := routes.Group("api/v1/admin")
	{
		//admin authentication
		admin.POST("/singup", controllers.AdminSignup)
		admin.POST("/login", controllers.AdminLogin)
		admin.Use(middleware.AuthMidleware())
		admin.Use(middleware.ProtectedRoutes())

		//user management
		admin.GET("/users", controllers.GetUsers)
		admin.GET("/users/user/:id", controllers.GetUser)
		admin.PATCH("/users/block/:id", controllers.BlockUser)
		admin.PATCH("/users/unblock/:id", controllers.UnblockUsers)

		//categories mangement
		admin.GET("/categories", controllers.GetCategories)
		admin.POST("/categories", controllers.AddCategories)
		admin.DELETE("/categories/:id", controllers.DeleteCategory)
		admin.PATCH("/categories/:id", controllers.EditCategory)

		//product management
		admin.GET("/products", controllers.GetProducts)
		admin.GET("/product/:id", controllers.GetProduct)
		admin.POST("/product", controllers.AddProduct)
		admin.PATCH("/product/:id", controllers.EditProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
		admin.GET("/products/filter", controllers.FilterProduct)

		//order managements
		admin.GET("/orders", controllers.OrdersList)
		admin.GET("/orders/:id/items", controllers.OrderItemList)
		admin.PUT("/orders/:id/status", controllers.ChangeOrderStaus)

	}

}

func UserRoutes(routes *gin.Engine) {
	user := routes.Group("api/v1/user")
	{
		//user signup
		user.POST("/signup", controllers.UserSignup)
		user.POST("/signup/resend-otp", controllers.ResendOtp)
		user.POST("signup/varify-otp", controllers.VerifyOtp)
		user.POST("/login", controllers.UserLogin)
		//user forgot password
		user.POST("/forgot-password/send-otp", controllers.PasswordSendOtp)
		user.POST("/forgot-password/resend-otp", controllers.ResendOtp)
		user.POST("/forgot-password/verify-otp", controllers.ResetPassOtpVerify)
		user.POST("/forgot-password/reset-password", controllers.ResetPassword)
		//google login
		user.GET("/google-login", controllers.GoogleLogin)
		user.GET("/google/callback", controllers.GoogleCallback)
		user.Use(middleware.AuthMidleware())
		//product listing
		user.GET("/search", controllers.SearchProduct)
		user.GET("/categories", controllers.GetCategories)
		user.GET("/products", controllers.GetProducts)
		user.GET("/product/:id", controllers.GetProduct)
		user.GET("/product/filter", controllers.FilterProduct)
		//User profile
		user.GET("/user-profile", controllers.GetUserProfile)
		user.PUT("/user-profile", controllers.EditUserProfile)
		user.PATCH("/user-profile/password-change", controllers.ChangePassword)
		//user address
		user.GET("/user-profile/addresses", controllers.GetAddress)
		user.POST("/user-profile/addresses", controllers.AddAddress)
		user.PUT("/user-profile/addresses/:id", controllers.EditAddress)
		user.DELETE("/user-profile/addresses/:id", controllers.DeleteAddress)
		//user cart
		user.POST("/cart/:product_id", controllers.AddToCart)
		user.DELETE("/cart/:product_id", controllers.CartRemoveItem)
		user.GET("/cart", controllers.GetCart)
		user.PUT("/cart/:id", controllers.EditCart)

		//order management
		user.POST("/order", controllers.PlaceOrder)
		user.GET("/order/:order_id", controllers.GetOrderDetails)
		user.GET("/user-profile/orders", controllers.ListOrders)
		user.PUT("/order/:order_id", controllers.CancelOrder)

	}

}
