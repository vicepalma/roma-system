package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

func RequireCoachOf(svc service.CoachService, param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		coachID := MustUserID(c) // <- usa el userID del contexto
		if coachID == "" {       // ya abortó con 401
			return
		}
		discipleID := c.Param(param)

		ok, err := svc.CanCoach(c.Request.Context(), coachID, discipleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
