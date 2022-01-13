package middleware

import (
	"net/http"
	"strings"
	"time"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token and verify token
		authorization := c.Request.Header["Authorization"]
		token := strings.Split(authorization[0], " ")[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "no token"})
			c.Abort()
		}

		claims, err := pkg.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "parse token error"})
			c.Abort()
		}
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token expired"})
			c.Abort()
		}

		c.Next()
	}
}
