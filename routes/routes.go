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
		// admin.Use(middleware.AuthMidleware())
		// admin.Use(middleware.ProtectedRoutes())

		//admin dashoboard
		admin.GET("/dashoboard/sales-details", controllers.GetSalesData)
		admin.GET("dashoboard/top-product", controllers.BestSellingProduct)
		admin.GET("/dashoboard/top-categories", controllers.BestCategories)
		admin.GET("/dashoboard/sales-graph", controllers.SalesGraph)

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

		//add offer
		admin.PUT("/product/offer/:product_id", controllers.AddProductOffter)
		admin.GET("/offers", controllers.ListOffer)
		admin.PUT("/offers/offer-remove/:product_id", controllers.RemoveOffer)

		//coupon mangement

		admin.POST("/coupon", controllers.CreateCoupan)
		admin.DELETE("/coupon/:id", controllers.DeleteCoupan)
		admin.PUT("/coupon/:id", controllers.EditCoupan)
		admin.GET("/coupons", controllers.ListCoupan)

		//order managements
		admin.GET("/orders", controllers.OrdersList)
		admin.GET("/orders/:id/items", controllers.OrderItemList)
		admin.PUT("/orders/:id/status", controllers.ChangeOrderStaus)

		//sales SalesReport
		admin.GET("/sales-report", controllers.SalesReport)
		admin.GET("/sales-report/pdf", controllers.SalesReportPDF)

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
		// user.Use(middleware.AuthMidleware())
		//product listing
		user.GET("/search", controllers.SearchProduct)
		user.GET("/categories", controllers.GetCategories)
		user.GET("/products", controllers.GetProducts)
		user.GET("/product/:id", controllers.GetProduct)
		user.GET("/product/filter", controllers.FilterProduct)
		//Review
		user.POST("/product/:id/review", middleware.AuthMidleware(), controllers.Addreview)
		user.PUT("/product/:id/review",middleware.AuthMidleware(),controllers.UpdateReview)
		user.GET("/product/:id/review", controllers.GetProductReview)
		user.DELETE("/product/:id/review", middleware.AuthMidleware(), controllers.DeleteReview)
		//User profile
		user.GET("/user-profile", middleware.AuthMidleware(), controllers.GetUserProfile)
		user.PUT("/user-profile", middleware.AuthMidleware(), controllers.EditUserProfile)
		user.PATCH("/user-profile/password-change", middleware.AuthMidleware(), controllers.ChangePassword)
		user.PUT("/user-profile/referral-bonus", middleware.AuthMidleware(), controllers.AddReferralPointsToWallet)
		user.GET("/user-profile/wallete-histories", middleware.AuthMidleware(), controllers.WalleteHistory)
		//user address
		user.GET("/user-profile/addresses", middleware.AuthMidleware(), controllers.GetAddress)
		user.POST("/user-profile/addresses", middleware.AuthMidleware(), controllers.AddAddress)
		user.PUT("/user-profile/addresses/:id", middleware.AuthMidleware(), controllers.EditAddress)
		user.DELETE("/user-profile/addresses/:id", middleware.AuthMidleware(), controllers.DeleteAddress)
		//user cart
		user.POST("/cart/:product_id", middleware.AuthMidleware(), controllers.AddToCart)
		user.DELETE("/cart/:product_id", middleware.AuthMidleware(), controllers.CartRemoveItem)
		user.GET("/cart", middleware.AuthMidleware(), controllers.GetCart)
		user.PUT("/cart/:id", middleware.AuthMidleware(), controllers.EditCart)

		//wishlistt history

		user.GET("/wishlist", middleware.AuthMidleware(), controllers.GetWishlist)
		user.POST("/wishlist/:product_id", middleware.AuthMidleware(), controllers.AddWishlistItems)
		user.DELETE("/wishlist/:id", middleware.AuthMidleware(), controllers.DeleteWishilistItems)

		//order management

		user.POST("/order", middleware.AuthMidleware(), controllers.PlaceOrder)
		user.GET("/order/:order_id", middleware.AuthMidleware(), controllers.GetOrderDetails)
		user.GET("/user-profile/orders", middleware.AuthMidleware(), controllers.ListOrders)

		user.PUT("/order/order-cancel/:id", middleware.AuthMidleware(), controllers.CancelOrder)
		user.PUT("/order/order-item/:id", middleware.AuthMidleware(), controllers.CancelOrderProduct)
		user.PUT("/order/:id/order-return", middleware.AuthMidleware(), controllers.OrderReturn)
		user.GET("/order/order-invoice/:id", middleware.AuthMidleware(), controllers.GenerateOrderInvoice)

		// //RazorPay payment
		user.GET("/order/payment", controllers.RenderRazorpay)
		user.POST("/order/payment", controllers.RazapayPayment)
		user.POST("/order/payment/verify-payment", controllers.VerifyPaymentHandler)
		user.POST("/order/payment/failed-payment",controllers.PaymentFailed)

	}

}
