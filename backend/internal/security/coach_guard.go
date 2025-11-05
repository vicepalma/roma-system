package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func RequireCoachOfQuery(svc service.CoachService, queryKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		coachID := UserID(c)
		if coachID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		discipleID := c.Query(queryKey)
		if _, err := uuid.Parse(discipleID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "invalid disciple_id"})
			c.Abort()
			return
		}
		ok, err := svc.CanCoach(c.Request.Context(), coachID, discipleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
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
