package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type MeHandler struct{ hist service.HistoryService }

func NewMeHandler(hist service.HistoryService) *MeHandler { return &MeHandler{hist: hist} }

// GET /api/me/today
func (h *MeHandler) GetToday(c *gin.Context) {
	userID := security.MustUserID(c)
	if userID == "" {
		return
	}

	out, err := h.hist.GetMeTodayFor(c, userID)
	if err != nil {
		// ErrNoDay -> 404 semántico
		if err.Error() == "no_day" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no_day", "detail": "no program day for today"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *MeHandler) GetTodayForDisciple(c *gin.Context) {
	discipleID := c.Param("id")

	out, err := h.hist.GetMeTodayFor(c, discipleID)
	if err != nil {
		if err.Error() == "no_day" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no_day"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
