package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type MeHandler struct {
	hist    service.HistoryService
	coach   service.CoachService
	session service.SessionService
}

func (h *MeHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/me")
	{
		g.GET("/assignment/active", h.getActiveAssignment)
		g.GET("/today", h.GetToday)
		g.GET("/today/disciple", h.GetTodayForDisciple)
		g.GET("/session/active", h.getActiveSession)
	}
}

func NewMeHandler(hist service.HistoryService, cs service.CoachService, ss service.SessionService) *MeHandler {
	return &MeHandler{hist: hist, coach: cs, session: ss}
}

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

func (h *MeHandler) getActiveAssignment(c *gin.Context) {
	uid := security.UserID(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	assign, err := h.coach.GetActiveAssignment(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	if assign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, assign)
}

func (h *MeHandler) getActiveSession(c *gin.Context) {
	uid := security.UserID(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	s, err := h.session.GetActiveOpenSessionForMe(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	if s == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, s)
}
