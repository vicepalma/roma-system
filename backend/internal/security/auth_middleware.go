package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userID"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_bearer"})
			return
		}
		tokStr := strings.TrimPrefix(h, "Bearer ")
		tok, claims, err := ParseAndValidate(tokStr)
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}
		if typ, _ := claims["typ"].(string); typ != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong_token_type"})
			return
		}
		sub, _ := claims["sub"].(string)
		if sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no_sub"})
			return
		}
		c.Set(CtxUserID, sub)
		c.Next()
	}
}
