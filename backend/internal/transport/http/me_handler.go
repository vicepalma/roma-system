package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type MeHandler struct{ hist service.HistoryService }

func NewMeHandler(hist service.HistoryService) *MeHandler { return &MeHandler{hist: hist} }

// GET /api/me/today
func (h *MeHandler) GetToday(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	discipleID := userID.(string)
	tz := c.DefaultQuery("tz", "America/Santiago")

	out, err := h.hist.GetMeTodayFor(c.Request.Context(), discipleID, tz)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no_day", "detail": "no assignment/day for today"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *MeHandler) GetTodayForDisciple(c *gin.Context) {
	discipleID := c.Param("id")
	tz := c.DefaultQuery("tz", "America/Santiago")

	out, err := h.hist.GetMeTodayFor(c.Request.Context(), discipleID, tz) // <-- usa HistoryService
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
