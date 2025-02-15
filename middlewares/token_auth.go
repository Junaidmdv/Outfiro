package middleware

import (
	"fmt"
	"net/http"
	"os"
	"outfiro/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		SingnedToken, err := jwt.ParseWithClaims(clientToken, &utils.JwtClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("TOKEN_SECRETE_KEY")), nil
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"code":   "StatusUnauthorized",
				"detail": "Invalid token ",
			})
			c.Abort()
			return
		}
		claim, ok := SingnedToken.Claims.(*utils.JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "invalid claims"})
			return
		}
		if claim.ExpiresAt != nil && claim.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    "StatusUnauthorized",
				"details": "user token is expired.Please login again",
			})
			return
		}
		fmt.Println("claims", claim)
		if claim.UserId == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		if claim.Email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		c.Set("user_id", claim.UserId)
		c.Set("user_email", claim.Email)
		c.Set("Issuer", claim.Issuer)
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
