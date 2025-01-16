package middleware

import (
	"errors"
	"fmt"
	"outfiro/utils"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header missing",
			})
			return
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid token formate",
			})
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return utils.Secretekey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"errors": fmt.Sprintf("Invalid token %v", err)})
			return
		}
		if Claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			if err := utils.ValidateToken(Claims); err != nil {
				c.JSON(401, gin.H{"error": fmt.Sprintf("invalid token %v", err)})
				c.Abort()
				return
			}
			c.Set("claims", Claims)
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"errors": "invalid token claims",
			})
			c.Abort()
			return
		}

	}
}
func AdminAuth()gin.HandlerFunc{
	return func(c *gin.Context){
		claims,exist:=c.Get("claims")
		if !exist{
			c.JSON(401,gin.H{"error":"token claims is missing"})
			c.Abort()
			return
		}
		ClaimsData,Ok:=claims.(utils.Claims)
		if !Ok {
			c.JSON(401,gin.H{"error":"invalid claims type"})
			c.Abort()
			return
		}
		if ClaimsData.UserRole !="admin"{
            c.JSON(401,gin.H{"error":"unotherised"})
			c.Abort()
			return
		}
      c.Next()
	}
}


