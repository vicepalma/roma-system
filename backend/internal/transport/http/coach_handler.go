package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type assignReq struct {
	DiscipleID string `json:"disciple_id" binding:"required"`
	ProgramID  string `json:"program_id"  binding:"required"`
	StartDate  string `json:"start_date"` // YYYY-MM-DD opcional
}

type CoachHandler struct {
	svc   service.CoachService
	hist  service.HistoryService
	users repository.UserRepository
}

type createInviteReq struct {
	Email string `json:"email" binding:"required,email"`
	TTLh  *int   `json:"ttl_h"` // opcional, default 168h (7 días)
}

type createLinkReq struct {
	DiscipleID string `json:"disciple_id"` // opcional para auto-vínculo
	AutoAccept bool   `json:"auto_accept"` // opcional (para pruebas)
}

func NewCoachHandler(svc service.CoachService, hist service.HistoryService, u repository.UserRepository) *CoachHandler {
	return &CoachHandler{svc: svc, hist: hist, users: u}
}

func (h *CoachHandler) Register(r *gin.RouterGroup) {
	grp := r.Group("/coach")
	{
		grp.POST("/links", h.createLink)
		grp.PATCH("/links/:id", h.updateLink)
		grp.GET("/links", h.listLinks)

		grp.GET("/disciples", h.listDisciples)
		grp.GET("/disciples/:id/today",
			security.RequireCoachOf(h.svc, "id"),
			h.getTodayForDisciple,
		)
		grp.POST("/assignments", h.assignProgram)
		grp.GET("/assignments", h.listAssignments)
		grp.GET("/disciples/:id/overview",
			security.RequireCoachOf(h.svc, "id"),
			h.getOverview,
		)
		grp.PATCH("/assignments/:id", h.patchAssignment)
		grp.GET("/assignments/:id/calendar", h.assignmentCalendar)
		// grp.POST("/invitations", security.RequireCoachOfSelf(h.svc), h.createInvite) // o AuthRequired si no tienes rol
		// grp.POST("/invitations/:code/accept", security.AuthRequired(), h.acceptInvite)
	}
}

// @Summary Crear vínculo maestro-discípulo
// @Description Crea invitación (o auto-vínculo si disciple_id vacío o igual al coach).
// @Tags coach
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body struct{DiscipleID string `json:"disciple_id"`; AutoAccept bool `json:"auto_accept"`} true "Payload"
// @Success 201 {object} domain.CoachLink
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/coach/links [post]

func (h *CoachHandler) createLink(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	coachID := userID.(string)

	var req createLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	discipleID := req.DiscipleID
	if discipleID == "" {
		discipleID = coachID // autovínculo
	}
	link, err := h.svc.CreateLink(c, coachID, discipleID, req.AutoAccept)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "create_failed", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, link)
}

type updateLinkReq struct {
	Action string `json:"action"` // "accept" | "reject"
}

func (h *CoachHandler) updateLink(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get(security.CtxUserID)
	actorID := userID.(string)

	var req updateLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	link, err := h.svc.UpdateLinkStatus(c, id, actorID, req.Action)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "update_failed", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

// @Summary Listar discípulos del coach
// @Tags coach
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string][]repository.DiscipleRow
// @Failure 401 {object} map[string]string
// @Router /api/coach/disciples [get]

func (h *CoachHandler) listLinks(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	uid := userID.(string)
	incoming, outgoing, err := h.svc.ListLinks(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list_failed", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"incoming": incoming, // invitaciones donde soy discípulo
		"outgoing": outgoing, // invitaciones donde soy coach
	})
}

func (h *CoachHandler) listDisciples(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	coachID := userID.(string)
	rows, err := h.svc.ListDisciples(c, coachID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list_failed", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows})
}

func (h *CoachHandler) assignProgram(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	coachID := userID.(string)

	var req assignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	var start time.Time
	var err error
	if req.StartDate != "" {
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": "invalid start_date (YYYY-MM-DD)"})
			return
		}
	}

	row, err := h.svc.AssignProgram(c, coachID, req.DiscipleID, req.ProgramID, start)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "assign_failed", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, row)
}

// @Summary Overview del discípulo (hoy + pivot + adherencia)
// @Description Requiere que el solicitante sea coach del discípulo.
// @Tags coach
// @Security BearerAuth
// @Produce json
// @Param id path string true "Disciple ID"
// @Param days query int false "Días (por defecto 14)"
// @Param metric query string false "volume|sets|reps (por defecto volume)"
// @Param tz query string false "IANA TZ (por defecto America/Santiago)"
// @Success 200 {object} struct{DiscipleID string `json:"disciple_id"`; MeToday interface{} `json:"me_today"`; Pivot service.PivotResponse `json:"pivot"`; Adherence struct{Days int `json:"days"`; DaysWithSets int `json:"days_with_sets"`; Rate float64 `json:"rate"`} `json:"adherence"`}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/coach/disciples/{id}/overview [get]

func (h *CoachHandler) getOverview(c *gin.Context) {
	discipleID := c.Param("id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	metric := c.DefaultQuery("metric", "volume")
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	coachID := security.MustUserID(c)
	out, err := h.svc.GetOverview(c.Request.Context(), coachID, discipleID, days, metric, tz)
	if err != nil {
		// Si adentro caímos por ErrNoDay, devolvé estructura vacía
		if errors.Is(err, service.ErrNoDay) {
			c.JSON(http.StatusOK, gin.H{
				"disciple_id": discipleID,
				"me_today": gin.H{
					"assignment_id":              nil,
					"day":                        nil,
					"prescriptions":              []domain.Prescription{},
					"current_session_id":         nil,
					"current_session_started_at": nil,
					"current_session_sets_count": 0,
				},
				"pivot": gin.H{
					"mode":    "by_exercise",
					"days":    days,
					"columns": []string{"date"},
					"rows":    []any{},
					"catalog": []any{},
				},
				"adherence": gin.H{
					"days":           days,
					"days_with_sets": 0,
					"rate":           0.0,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *CoachHandler) listAssignments(c *gin.Context) {
	userID, _ := c.Get(security.CtxUserID)
	coachID := userID.(string)

	discipleID := c.Query("disciple_id")
	var disciplePtr *string
	if discipleID != "" {
		disciplePtr = &discipleID
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.svc.ListAssignments(c, coachID, disciplePtr, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list_failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CoachHandler) getTodayForDisciple(c *gin.Context) {
	ctx := c.Request.Context()
	discipleID := c.Param("id")

	// Autorización (como ya la tienes)
	userID, _ := c.Get(security.CtxUserID)
	coachID := userID.(string)

	// Autorización: coach válido para ese discípulo
	ok, err := h.svc.CanCoach(ctx, coachID, discipleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// TZ
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	me, err := h.hist.GetMeTodayFor(ctx, discipleID, tz)
	if err != nil {
		// Si no hay "hoy", devolvemos 200 con shape vacío (para que el front no caiga)
		if errors.Is(err, service.ErrNoDay) || errors.Is(err, repository.ErrNoDay) {
			c.JSON(http.StatusOK, gin.H{
				"assignment_id":              nil,
				"day":                        nil,
				"prescriptions":              []repository.MeTodayPrescription{},
				"current_session_id":         nil,
				"current_session_started_at": nil,
				"current_session_sets_count": nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}

	// ¡OJO!: NO desreferenciar punteros ni llamar métodos sobre punteros opcionales.
	// Pasamos todo tal cual; gin/json serializa `nil` como `null`.

	// day: mantener puntero (null-safe)
	var day any = nil
	if me != nil && me.Day != nil {
		day = me.Day
	}

	// prescripciones: lista vacía si viene nil
	presc := []repository.MeTodayPrescription{}
	if me != nil && me.Prescriptions != nil {
		presc = me.Prescriptions
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment_id":              me.AssignmentID,
		"day":                        day,
		"prescriptions":              presc,
		"current_session_id":         me.CurrentSessionID,
		"current_session_started_at": me.CurrentSessionStartedAt, // ← puntero tal cual (puede ser nil)
		"current_session_sets_count": me.CurrentSessionSetsCount, // ← puntero tal cual (puede ser nil)
	})
}

func (h *CoachHandler) patchAssignment(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		EndDate  *string `json:"end_date"`  // "YYYY-MM-DD"
		IsActive *bool   `json:"is_active"` // true/false
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	var endPtr *time.Time
	if body.EndDate != nil && *body.EndDate != "" {
		d, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_end_date"})
			return
		}
		endPtr = &d
	}
	asg, err := h.svc.UpdateAssignment(c.Request.Context(), id, endPtr, body.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_update"})
		return
	}
	c.JSON(http.StatusOK, asg)
}

func (h *CoachHandler) assignmentCalendar(c *gin.Context) {
	id := c.Param("id")
	fromStr := c.DefaultQuery("from", time.Now().Format("2006-01-02"))
	toStr := c.DefaultQuery("to", time.Now().AddDate(0, 0, 13).Format("2006-01-02"))

	from, err1 := time.Parse("2006-01-02", fromStr)
	to, err2 := time.Parse("2006-01-02", toStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_range"})
		return
	}

	items, err := h.svc.AssignmentCalendar(c.Request.Context(), id, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_build_calendar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "from": fromStr, "to": toStr})
}
