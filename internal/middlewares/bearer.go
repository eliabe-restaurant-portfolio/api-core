package middlewares

import (
	"net/http"
	"strings"

	"github.com/eliabe-portfolio/restaurant-app/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type SubPayload struct {
	UUID string `json:"uuid"`
}

func (m middlewares) BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		publicPem, err := jwt.LoadPublicKeyFromFile("./storage/public_key.pem")
		if err != nil {
			c.Abort()
			return
		}

		claims, err := jwt.Read(jwt.JwtReadInput{
			ExternalToken: tokenString,
			PublicPem:     []byte(publicPem),
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload: sub claim missing or not string"})
			c.Abort()
			return
		}

		c.Set("actor", sub)
		c.Next()
	}
}
