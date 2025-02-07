package middleware

import (
	"outfiro/utils"
	"strings"
	"github.com/gin-gonic/gin"
)

func AuthMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ExtractedToken := c.Request.Header.Get("Authorization")
		if ExtractedToken == " " {
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized(401),",
				"message": "no authentication header is provided",
			})
			c.Abort()
			return
		}
		token := strings.Split(ExtractedToken, "Bearer ")
		if len(token) != 2 {
			c.JSON(401, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized",
				"message": "No authentication header is provided",
			})
			c.Abort()
			return
		}
		clientToken := strings.TrimSpace(token[1])
		_, err := utils.ValidateToken(clientToken, c)
		if err != nil {
			c.JSON(403, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized",
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
func ProtectedRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		Isuuer, exist := c.Get("Issuer")
		if !exist {
			c.JSON(403, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized(403)",
				"message": "Missing Issuer data",
			})
			c.Abort()
			return
		}
		if Isuuer != "admin" {
			c.JSON(403, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized(403)",
				"message": "only admin can access this routes",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}
