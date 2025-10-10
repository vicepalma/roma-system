package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type SessionHandler struct{ svc service.SessionService }

func NewSessionHandler(s service.SessionService) *SessionHandler { return &SessionHandler{svc: s} }

func (h *SessionHandler) Register(r *gin.RouterGroup) {
	r.POST("/sessions", h.start)                // crea sesión
	r.GET("/sessions/:id", h.get)               // detalle + sets + cardio
	r.POST("/sessions/:id/sets", h.addSet)      // agrega set
	r.PATCH("/sessions/:id", h.update)          // notas/fecha
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

func (h *SessionHandler) update(c *gin.Context) {
	id := c.Param("id")
	type req struct {
		PerformedAt *string `json:"performed_at"`
		Notes       *string `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
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
	if err := h.svc.Update(c, uid(c), id, ts, body.Notes); err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.Status(204)
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
