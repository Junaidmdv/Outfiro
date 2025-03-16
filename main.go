package main

import (
	"os"
	"outfiro/database"
	"outfiro/models"
	"outfiro/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("templates/app.html")
	
	database.LoadEnv()

	database.DbConnection()
	models.Migrate()
	
	routes.AdminRoutes(router)
	routes.UserRoutes(router)


	router.Run(os.Getenv("PORT"))
}
