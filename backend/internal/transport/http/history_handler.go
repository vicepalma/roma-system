package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type HistoryHandler struct {
	svc   service.HistoryService
	defTz string
}

func NewHistoryHandler(s service.HistoryService, defaultTZ string) *HistoryHandler {
	if defaultTZ == "" {
		defaultTZ = "UTC"
	}
	return &HistoryHandler{svc: s, defTz: defaultTZ}
}

func (h *HistoryHandler) Register(r *gin.RouterGroup) {
	grp := r.Group("/history")
	{
		grp.GET("", h.history)         // sesiones recientes  /api/history?disciple_id=&from=&to=&group=day|session&tz=&limit=&offset=
		grp.GET("/summary", h.summary) // nuevos endpoints de resumen
		grp.GET("/summary/pivot", h.summaryPivot)
		grp.GET("/prs", h.prs)

		// /api/disciples/:id/sessions?from=&to=&tz=&limit=&offset=
		grp.GET("/disciples/:id/sessions", h.sessions)

		// /api/disciples/:id/days?from=&to=&tz=&limit=&offset=
		grp.GET("/disciples/:id/days", h.planVsDone)
	}
}

func (h *HistoryHandler) history(c *gin.Context) {
	discipleID := c.Query("disciple_id")
	if discipleID == "" {
		// fallback: si no viene, usa el usuario autenticado
		discipleID = security.UserID(c)
	}
	if discipleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "disciple_id_required"})
		return
	}
	group := c.DefaultQuery("group", "session")
	tz := c.DefaultQuery("tz", h.defTz)

	from, to := parseDates(c.Query("from")), parseDates(c.Query("to"))
	limit, offset := parsePag(c.DefaultQuery("limit", "50")), parsePag(c.DefaultQuery("offset", "0"))

	data, total, err := h.svc.History(c.Request.Context(), discipleID, tz, group, from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data, "total": total, "limit": limit, "offset": offset})
}

func (h *HistoryHandler) sessions(c *gin.Context) {
	discipleID := c.Param("id")
	tz := c.DefaultQuery("tz", h.defTz)
	from, to := parseDates(c.Query("from")), parseDates(c.Query("to"))
	limit, offset := parsePag(c.DefaultQuery("limit", "50")), parsePag(c.DefaultQuery("offset", "0"))

	data, total, err := h.svc.ListSessions(c.Request.Context(), discipleID, tz, from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data, "total": total, "limit": limit, "offset": offset})
}

func (h *HistoryHandler) planVsDone(c *gin.Context) {
	discipleID := c.Param("id")
	tz := c.DefaultQuery("tz", h.defTz)
	from, to := parseDates(c.Query("from")), parseDates(c.Query("to"))
	limit, offset := parsePag(c.DefaultQuery("limit", "50")), parsePag(c.DefaultQuery("offset", "0"))

	data, total, err := h.svc.PlanVsDone(c.Request.Context(), discipleID, tz, from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data, "total": total, "limit": limit, "offset": offset})
}

func (h *HistoryHandler) summary(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	mode := strings.ToLower(c.DefaultQuery("mode", "by_exercise"))
	tz := c.DefaultQuery("tz", "America/Santiago")

	userID, _ := c.Get(security.CtxUserID)
	uid := userID.(string)

	switch mode {
	case "by_muscle":
		rows, err := h.svc.GetDailyByMuscle(c, uid, days, tz)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"mode": "by_muscle", "items": rows, "days": clamp(days)})
	default:
		includeCatalog := strings.EqualFold(c.DefaultQuery("include", ""), "catalog")
		rows, catalog, err := h.svc.GetDailyByExercise(c, uid, days, includeCatalog, tz)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
			return
		}
		resp := gin.H{"mode": "by_exercise", "items": rows, "days": clamp(days)}
		if includeCatalog {
			resp["catalog"] = catalog
		}
		c.JSON(http.StatusOK, resp)
	}
}

func (h *HistoryHandler) prs(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	rows, err := h.svc.GetPRs(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows})
}

func clamp(n int) int {
	if n <= 0 || n > 180 {
		return 14
	}
	return n
}

func (h *HistoryHandler) summaryPivot(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	mode := strings.ToLower(c.DefaultQuery("mode", "by_exercise"))
	includeCatalog := strings.EqualFold(c.DefaultQuery("include", ""), "catalog")
	metric := c.DefaultQuery("metric", "volume")
	tz := c.DefaultQuery("tz", "America/Santiago")

	userID, _ := c.Get(security.CtxUserID)
	uid := userID.(string)

	switch mode {
	case "by_muscle":
		resp, err := h.svc.GetPivotByMuscle(c, uid, days, metric, tz)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	default:
		resp, err := h.svc.GetPivotByExercise(c, uid, days, includeCatalog, metric, tz)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// helpers
func parseDates(s string) *time.Time {
	if s == "" {
		return nil
	}
	// acepta YYYY-MM-DD o RFC3339
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
	}
	return nil
}
func parsePag(s string) int {
	n, _ := strconv.Atoi(s)
	if n < 0 {
		n = 0
	}
	return n
}
