package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"workerunion/db/handlers"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token and verify token
		authorization := c.Request.Header["Authorization"]
		token := strings.Split(authorization[0], " ")[1]
		fmt.Println("----token--", token)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no token"})
			return
		}

		claims, err := pkg.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "parse token error"})
			return
		}
		fmt.Println("-----claims err", err, err != nil, 1)
		if claims.ExpiresAt < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token expired"})
			return
		}

		fmt.Println("login--", claims.UserId)
		query := map[string]interface{}{
			"id": claims.UserId,
		}
		users := handlers.FindUsers(query)
		if len(users) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no account find"})
			return
		}
		user := users[0]
		c.Set("currentUser", user)

		c.Next()
	}
}
