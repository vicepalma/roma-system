package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

func CoachGuard(svc service.CoachService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(CtxUserID)
		coachID := userID.(string)
		discipleID := c.Param("disciple_id")
		if discipleID == "" {
			discipleID = c.Query("disciple_id")
		}
		ok, err := svc.CanCoach(c, coachID, discipleID)
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
