package routes

import (
	
    "outfiro/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.Engine) {
	admin := routes.Group("/admin")
	{
		admin.POST("singup", controllers.AdminSignup)
		admin.POST("login", controllers.AdminLogin)
	
		admin.GET("/users", controllers.GetUsers)
		admin.GET("/users/user/:id",controllers.GetUser)
		admin.PATCH("amdin/users/block/:id", controllers.BlockUser)
		admin.PATCH("/admin/users/ublock/:id", controllers.UnblockUsers)
		admin.GET("/categories", controllers.GetCategories)
		admin.POST("/categories", controllers.AddCategories)
		admin.DELETE("/categories", controllers.DeleteCategory)
		admin.PATCH("/categories", controllers.EditCategory)
		admin.GET("/products", controllers.GetProducts)
		admin.GET("/product/:id", controllers.GetProduct)
		admin.POST("/product", controllers.AddProduct)
		admin.PUT("/product/id", controllers.EditProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
		// admin.POST("/products/:id/images", controllers.AddProductImages)
	}

}

func UserRoutes(routes *gin.Engine) {
	user := routes.Group("/user")
	{
		user.POST("signup", controllers.UserSignup)
		user.POST("sigup", controllers.UserLogin)
		user.POST("signup/varify-otp", controllers.VerifyOtp)
		user.POST("login", controllers.UserLogin)
		user.GET("/search", controllers.SearchProduct)
		user.GET("/categories", controllers.GetCategories)
		user.GET("/products", controllers.GetProducts)
		user.GET("/product/:id", controllers.GetProduct)
		
	}

}
