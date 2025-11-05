// internal/security/auth_middleware.go
package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const CtxUserID = "user_id"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Deja pasar preflight
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		prefix := "bearer "
		if len(auth) < 8 || strings.ToLower(auth[:len(prefix)]) != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		raw := strings.TrimSpace(auth[len(prefix):])
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Parse/valida usando tu helper existente
		token, claims, err := ParseAndValidate(raw)
		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// (Opcional pero recomendado) clock skew pequeño
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().After(time.Unix(int64(exp), 0).Add(15 * time.Second)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		}

		// Exigir typ=access si está presente
		if typ, _ := claims["typ"].(string); typ != "" && typ != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		sub, _ := claims["sub"].(string)
		if sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Listo: deja el user en el contexto para guards/handlers
		c.Set(CtxUserID, sub)
		c.Next()
	}
}

// Si usas RegisteredClaims en otros lados, helper para revisar exp/nbf opcional:
func validateStd(claims *jwt.RegisteredClaims) bool {
	now := time.Now()
	// 15s skew
	skew := 15 * time.Second
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time.Add(skew)) {
		return false
	}
	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time.Add(-skew)) {
		return false
	}
	return true
}
