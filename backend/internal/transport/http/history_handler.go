package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type HistoryHandler struct{ svc service.HistoryService }

func NewHistoryHandler(s service.HistoryService) *HistoryHandler { return &HistoryHandler{svc: s} }

func (h *HistoryHandler) Register(r *gin.RouterGroup) {
	r.GET("/history", h.history)         // sesiones recientes
	r.GET("/history/summary", h.summary) // nuevos endpoints de resumen
	r.GET("/history/summary/pivot", h.summaryPivot)
	r.GET("/prs", h.prs)
}

func (h *HistoryHandler) history(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	userID, _ := c.Get(security.CtxUserID)
	data, err := h.svc.GetHistory(c, userID.(string), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
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
