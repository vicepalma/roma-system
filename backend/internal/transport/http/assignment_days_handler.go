package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type AssignmentDaysHandler struct {
	svc service.AssignmentDaysService
}

func NewAssignmentDaysHandler(s service.AssignmentDaysService) *AssignmentDaysHandler {
	return &AssignmentDaysHandler{svc: s}
}

func (h *AssignmentDaysHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/assignments")
	{
		g.GET(":assignmentId/days", h.list)
	}
}

func (h *AssignmentDaysHandler) list(c *gin.Context) {
	assignmentID := c.Param("assignmentId")
	if _, err := uuid.Parse(assignmentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "invalid assignmentId"})
		return
	}
	reqID := security.UserID(c)
	items, err := h.svc.List(c.Request.Context(), reqID, assignmentID)
	if err != nil {
		if err.Error() == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": []any{}}) // también podrías devolver 200 vacío si prefieres
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
