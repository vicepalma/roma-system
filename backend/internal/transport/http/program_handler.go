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
	g := r.Group("/programs")
	{
		g.GET("", h.listMine)             // GET /programs
		g.POST("", h.createProgram)       // POST /programs
		g.GET("/:id", h.get)              // GET /programs/:id
		g.PUT("/:id", h.update)           // PUT /programs/:id
		g.DELETE("/:id", h.delete)        // DELETE /programs/:id
		g.POST("/:id/version", h.version) // POST /programs/:id/version
		g.GET("/:id/versions", h.versions)

		g.POST("/:id/weeks", h.addWeek)
		g.GET("/:id/weeks", h.listWeeks)
		g.POST("/:id/weeks/:weekId/days", h.addDay)
		g.GET("/:id/weeks/:weekId/days", h.listDays)
		g.PUT("/:id/weeks/:weekId/days/:dayId", h.updateDay)
		g.DELETE("/:id/weeks/:weekId/days/:dayId", h.deleteDay)

		g.GET("/programs/:id/weeks/:weekId/days", h.listDays)
		g.PUT("/programs/:id/weeks/:weekId/days/:dayId", h.updateDay)
		g.DELETE("/programs/:id/weeks/:weekId/days/:dayId", h.deleteDay)

		// Prescripciones
		g.GET("/days/:dayId/prescriptions", h.listPresc)
		g.POST("/days/:dayId/prescriptions", h.addPrescription)
		g.PUT("/prescriptions/:id", h.updatePresc)
		g.DELETE("/prescriptions/:id", h.deletePresc)
		g.PATCH("/prescriptions/reorder", h.reorderPresc)
		g.DELETE("/:id/weeks/:weekId", h.deleteWeek)
	}

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

func (h *ProgramHandler) get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProgramHandler) update(c *gin.Context) {
	id := c.Param("id")
	var in service.UpdateProgram
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProgramHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProgramHandler) version(c *gin.Context) {
	id := c.Param("id")
	v, err := h.svc.NewVersion(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_version"})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *ProgramHandler) versions(c *gin.Context) {
	id := c.Param("id")
	items, err := h.svc.ListVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ProgramHandler) listWeeks(c *gin.Context) {
	programID := c.Param("id")
	items, err := h.svc.ListWeeks(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ProgramHandler) deleteWeek(c *gin.Context) {
	programID := c.Param("id")
	weekID := c.Param("weekId")

	if programID == "" || weekID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_params"})
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.DeleteWeek(ctx, programID, weekID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProgramHandler) listDays(c *gin.Context) {
	weekID := c.Param("weekId")
	items, err := h.svc.ListDays(c.Request.Context(), weekID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
func (h *ProgramHandler) updateDay(c *gin.Context) {
	type body struct {
		DayIndex *int    `json:"day_index"`
		Notes    *string `json:"notes"`
	}
	var b body
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	dayID := c.Param("dayId")
	d, err := h.svc.UpdateDay(c.Request.Context(), dayID, b.Notes, b.DayIndex)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, d)
}
func (h *ProgramHandler) deleteDay(c *gin.Context) {
	dayID := c.Param("dayId")
	if err := h.svc.DeleteDay(c.Request.Context(), dayID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Prescriptions
func (h *ProgramHandler) listPresc(c *gin.Context) {
	dayID := c.Param("dayId")
	items, err := h.svc.ListPrescriptions(c.Request.Context(), dayID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ProgramHandler) updatePresc(c *gin.Context) {
	id := c.Param("id")
	var in service.UpdatePrescription
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	pr, err := h.svc.UpdatePrescription(c.Request.Context(), id, in)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, pr)
}
func (h *ProgramHandler) deletePresc(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePrescription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *ProgramHandler) reorderPresc(c *gin.Context) {
	type body struct {
		DayID string   `json:"day_id" binding:"required"`
		Order []string `json:"order" binding:"required,min=1"` // array de prescription IDs en orden
	}
	var b body
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.svc.ReorderPrescriptions(c.Request.Context(), b.DayID, b.Order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_reorder"})
		return
	}
	c.Status(http.StatusNoContent)
}

// (opcional) parse helpers si las necesitas
func parseInt(q string, def int) int {
	if q == "" {
		return def
	}
	v, err := strconv.Atoi(q)
	if err != nil {
		return def
	}
	return v
}

func parseDate(q string, def time.Time) time.Time {
	if q == "" {
		return def
	}
	t, err := time.Parse("2006-01-02", q)
	if err != nil {
		return def
	}
	return t
}
