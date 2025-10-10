package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type ProgramHandler struct{ svc service.ProgramService }

func NewProgramHandler(s service.ProgramService) *ProgramHandler { return &ProgramHandler{svc: s} }

func (h *ProgramHandler) Register(r *gin.RouterGroup) {
	r.POST("/programs", h.createProgram)
	r.GET("/programs", h.listMine)

	r.POST("/programs/:id/weeks", h.addWeek)
	r.POST("/weeks/:id/days", h.addDay)
	r.POST("/days/:id/prescriptions", h.addPrescription)

	r.POST("/assignments", h.assign)
	r.GET("/me/today", h.meToday)
}

func userID(c *gin.Context) string {
	v, _ := c.Get(security.CtxUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (h *ProgramHandler) createProgram(c *gin.Context) {
	type req struct {
		Title string  `json:"title" binding:"required,min=2"`
		Notes *string `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	p, err := h.svc.CreateProgram(c, userID(c), body.Title, body.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProgramHandler) listMine(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.svc.ListMyPrograms(c, userID(c), limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error"})
		return
	}
	c.JSON(200, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

func (h *ProgramHandler) addWeek(c *gin.Context) {
	id := c.Param("id")
	type req struct {
		WeekIndex int `json:"week_index" binding:"required,min=1"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
		return
	}
	w, err := h.svc.AddWeek(c, id, body.WeekIndex)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, w)
}

func (h *ProgramHandler) addDay(c *gin.Context) {
	id := c.Param("id")
	type req struct {
		DayIndex int     `json:"day_index" binding:"required,min=1"`
		Notes    *string `json:"notes"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
		return
	}
	d, err := h.svc.AddDay(c, id, body.DayIndex, body.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, d)
}

func (h *ProgramHandler) addPrescription(c *gin.Context) {
	id := c.Param("id")
	var body domain.Prescription
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	body.DayID = id
	p, err := h.svc.AddPrescription(c, &body)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, p)
}

func (h *ProgramHandler) assign(c *gin.Context) {
	type req struct {
		ProgramID  string  `json:"program_id" binding:"required"`
		DiscipleID string  `json:"disciple_id" binding:"required"`
		StartDate  string  `json:"start_date" binding:"required"` // YYYY-MM-DD
		EndDate    *string `json:"end_date"`
	}
	var body req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
		return
	}
	start, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad_date"})
		return
	}
	var endPtr *time.Time
	if body.EndDate != nil && *body.EndDate != "" {
		e, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad_date_end"})
			return
		}
		endPtr = &e
	}
	a, err := h.svc.Assign(c, body.ProgramID, body.DiscipleID, userID(c), start, endPtr)
	if err != nil {
		c.JSON(500, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	c.JSON(201, a)
}

func (h *ProgramHandler) meToday(c *gin.Context) {
	day, prescs, err := h.svc.MyToday(c, userID(c), time.Now().UTC())
	if err != nil {
		c.JSON(404, gin.H{"error": "no_day", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"day": day, "prescriptions": prescs})
}
