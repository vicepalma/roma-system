// internal/security/coach_guard_alias.go
package security

import (
	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

// Atajo para rutas con param ":id" (discipleID) usando el mismo guard base.
func RequireCoachOfSelf(svc service.CoachService) gin.HandlerFunc {
	return RequireCoachOf(svc, "id")
}
