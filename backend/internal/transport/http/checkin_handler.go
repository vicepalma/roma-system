package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
	"gorm.io/gorm"
)

type CheckinHandler struct {
	svc service.CheckinService
	db  *gorm.DB
}

func NewCheckinHandler(svc service.CheckinService, db *gorm.DB) *CheckinHandler {
	return &CheckinHandler{svc: svc, db: db}
}

func (h *CheckinHandler) Register(r *gin.RouterGroup) {
	r.POST("/checkins", security.RequireRole(h.db, "disciple"), h.create)
	r.GET("/checkins", security.RequireRole(h.db, "disciple"), h.listMine)
	r.GET("/checkins/:id", h.get)
	r.GET("/coach/disciples/:id/checkins", security.RequireRole(h.db, "coach"), h.listForCoach)
}

func (h *CheckinHandler) create(c *gin.Context) {
	type req struct {
		CheckedAt string   `json:"checked_at"`
		WeightKG  *float64 `json:"weight_kg"`
		Notes     *string  `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	checkedAt, ok := parseCheckinDate(c, body.CheckedAt)
	if !ok {
		return
	}
	if body.WeightKG != nil && *body.WeightKG <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_weight"})
		return
	}
	checkin, err := h.svc.Create(c.Request.Context(), security.UserID(c), checkedAt, body.WeightKG, cleanOptionalText(body.Notes))
	if err != nil {
		if errors.Is(err, service.ErrInvalidCheckin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusCreated, checkin)
}

func (h *CheckinHandler) listMine(c *gin.Context) {
	h.listByDisciple(c, security.UserID(c))
}

func (h *CheckinHandler) listForCoach(c *gin.Context) {
	discipleID := c.Param("id")
	ok, err := security.CanAccessDisciple(h.db.WithContext(c.Request.Context()), security.UserID(c), discipleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	h.listByDisciple(c, discipleID)
}

func (h *CheckinHandler) listByDisciple(c *gin.Context, discipleID string) {
	limit, offset := parseCheckinPag(c.DefaultQuery("limit", "50")), parseCheckinPag(c.DefaultQuery("offset", "0"))
	items, total, err := h.svc.List(c.Request.Context(), discipleID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

func (h *CheckinHandler) get(c *gin.Context) {
	checkin, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	ok, err := security.CanAccessDisciple(h.db.WithContext(c.Request.Context()), security.UserID(c), checkin.DiscipleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, checkin)
}

func parseCheckinDate(c *gin.Context, raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), true
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_checked_at"})
		return time.Time{}, false
	}
	return parsed, true
}

func cleanOptionalText(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func parseCheckinPag(s string) int {
	n, _ := strconv.Atoi(s)
	if n < 0 {
		return 0
	}
	return n
}
