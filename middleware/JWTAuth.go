package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

		header := c.Request.Header.Get("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "forbidden", "message": "missing access token"})
			return
		}

		var tokenString string = strings.ReplaceAll(header, "Bearer ", "")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return hmacSampleSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userId", claims["userId"])
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "forbidden", "message": err.Error()})
			return
		}
		c.Next()
	}
}
