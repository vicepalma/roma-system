package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MustUserID(c *gin.Context) string {
	val, ok := c.Get(CtxUserID) // CtxUserID = "userID" en tu AuthRequired
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return ""
	}
	id, _ := val.(string)
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_user"})
		c.Abort()
		return ""
	}
	return id
}
