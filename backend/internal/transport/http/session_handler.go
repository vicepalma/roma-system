package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type SessionHandler struct{ svc service.SessionService }

func NewSessionHandler(s service.SessionService) *SessionHandler { return &SessionHandler{svc: s} }

func (h *SessionHandler) Register(r *gin.RouterGroup) {
	r.POST("/sessions", h.start)           // crea sesión
	r.GET("/sessions/:id", h.get)          // detalle + sets + cardio
	r.POST("/sessions/:id/sets", h.addSet) // agrega set
	r.GET("/sessions/:id/sets", h.GetSets)
	r.DELETE("/sessions/:id/sets/:setId", h.deleteSet)
	r.PATCH("/sessions/:id", h.patchSession)    // notas/fecha
	r.POST("/sessions/:id/cardio", h.addCardio) // agrega cardio
}

func uid(c *gin.Context) string {
	v, _ := c.Get(security.CtxUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (h *SessionHandler) start(c *gin.Context) {
	type req struct {
		AssignmentID string  `json:"assignment_id" binding:"required"`
		DayID        string  `json:"day_id" binding:"required"`
		PerformedAt  *string `json:"performed_at"` // RFC3339 opcional
		Notes        *string `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	var ts *time.Time
	if body.PerformedAt != nil && *body.PerformedAt != "" {
		t, err := time.Parse(time.RFC3339, *body.PerformedAt)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad_datetime"})
			return
		}
		ts = &t
	}
	sess, err := h.svc.Start(c, uid(c), body.AssignmentID, body.DayID, ts, body.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, sess)
}

func (h *SessionHandler) get(c *gin.Context) {
	id := c.Param("id")
	sess, sets, cardio, err := h.svc.Get(c, uid(c), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not_found"})
		return
	}
	c.JSON(200, gin.H{"session": sess, "sets": sets, "cardio": cardio})
}

func (h *SessionHandler) addSet(c *gin.Context) {
	id := c.Param("id")
	type req struct {
		PrescriptionID string   `json:"prescription_id" binding:"required"`
		SetIndex       int      `json:"set_index" binding:"required,min=1"`
		Weight         *float64 `json:"weight"`
		Reps           int      `json:"reps" binding:"required,min=0"`
		RPE            *float32 `json:"rpe"`
		ToFailure      bool     `json:"to_failure"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	row, err := h.svc.AddSet(c, uid(c), id, body.PrescriptionID, body.SetIndex, body.Weight, body.Reps, body.RPE, body.ToFailure)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, row)
}

func (h *SessionHandler) addCardio(c *gin.Context) {
	id := c.Param("id")
	type req struct {
		Modality string  `json:"modality" binding:"required"`
		Minutes  int     `json:"minutes" binding:"required,min=1"`
		HRMin    *int    `json:"target_hr_min"`
		HRMax    *int    `json:"target_hr_max"`
		Notes    *string `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
		return
	}
	seg, err := h.svc.AddCardio(c, uid(c), id, body.Modality, body.Minutes, body.HRMin, body.HRMax, body.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, seg)
}

func (h *SessionHandler) GetSets(c *gin.Context) {
	// Auth: misma protección que POST /sessions/:id/sets
	userID, _ := c.Get(security.CtxUserID)
	actorID := userID.(string)

	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "missing session id"})
		return
	}

	var prescPtr *string
	if v := c.Query("prescription_id"); v != "" {
		prescPtr = &v
	}

	limit := 0
	offset := 0
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "invalid limit"})
			return
		}
		limit = n
	}
	if v := c.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "invalid offset"})
			return
		}
		offset = n
	}

	items, total, err := h.svc.ListSets(c.Request.Context(), actorID, sessionID, prescPtr, limit, offset)
	if err != nil {
		switch {
		case err.Error() == "forbidden":
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			// Si quieres distinguir 404, puedes mapear err == gorm.ErrRecordNotFound a 404.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Respuesta
	resp := gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SessionHandler) deleteSet(c *gin.Context) {
	_ = c.Param("id")
	setID := c.Param("setId")
	if err := h.svc.DeleteSet(c.Request.Context(), setID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SessionHandler) patchSession(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		PerformedAt *string `json:"performed_at"` // ISO8601
		Notes       *string `json:"notes"`
		Status      *string `json:"status"`   // open|closed
		EndedAt     *string `json:"ended_at"` // ISO8601
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	var tptr *time.Time
	if body.PerformedAt != nil && *body.PerformedAt != "" {
		t, err := time.Parse(time.RFC3339, *body.PerformedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_performed_at"})
			return
		}
		tptr = &t
	}

	var tend *time.Time
	if body.EndedAt != nil && *body.EndedAt != "" {
		t, err := time.Parse(time.RFC3339, *body.EndedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_ended_at"})
			return
		}
		tend = &t
	}

	out, err := h.svc.PatchSession(c.Request.Context(), id, tptr, body.Notes, body.Status, tend)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
