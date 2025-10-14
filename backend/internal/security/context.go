// internal/security/context.go
package security

import "github.com/gin-gonic/gin"

const CtxUserIDKey = "user_id" // debe coincidir con el c.Set(...) del AuthRequired

// UserID devuelve el ID de usuario (sub) puesto por el middleware de auth en el contexto.
// Si no existe, retorna "".
func UserID(c *gin.Context) string {
	v, ok := c.Get(CtxUserIDKey)
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
