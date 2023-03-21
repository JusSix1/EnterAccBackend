package middlewares

import (
	"net/http"
	"strings"

	"github.com/JusSix1/TwitterAccountDataBase/service"
	"github.com/gin-gonic/gin"
)

// validates token // เพื่อยืนยันว่ามี token จริง
func AuthorizesUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No Authorization for User header provided"})
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect Format of Authorization User Token"})
			return
		}

		jwtWrapper := service.JwtWrapperUser{
			SecretKey: "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
			Issuer:    "AuthService",
		}

		// ยืนยัน
		_, err := jwtWrapper.ValidateTokenUser(clientToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return

		}
		c.Next()
	}

}

func AuthorizesAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No Authorization for Admin header provided"})
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect Format of Authorization Admin Token"})
			return
		}

		jwtWrapper := service.JwtWrapperAdmin{
			SecretKey: "QCaxlcmpCdvosopjvpNKvlkdnvihiwcCAC",
			Issuer:    "AuthService",
		}

		// ยืนยัน
		_, err := jwtWrapper.ValidateTokenAdmin(clientToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return

		}
		c.Next()
	}

}
