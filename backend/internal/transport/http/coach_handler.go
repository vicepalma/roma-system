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

type assignReq struct {
	DiscipleID string `json:"disciple_id" binding:"required"`
	ProgramID  string `json:"program_id"  binding:"required"`
	StartDate  string `json:"start_date"` // YYYY-MM-DD opcional
}

type CoachHandler struct {
	svc  service.CoachService
	hist service.HistoryService
}

func NewCoachHandler(s service.CoachService) *CoachHandler { return &CoachHandler{svc: s} }

func (h *CoachHandler) Register(r *gin.RouterGroup) {
	grp := r.Group("/coach")
	{
		grp.POST("/links", h.createLink)
		grp.PATCH("/links/:id", h.updateLink)
		grp.GET("/links", h.listLinks)

		grp.GET("/disciples", h.listDisciples)
		grp.GET("/disciples/:id/today",
			security.RequireCoachOf(h.svc, "id"),
			h.GetTodayForDisciple,
		)
		grp.POST("/assignments", h.assignProgram)
		grp.GET("/assignments", h.listAssignments)
		grp.GET("/disciples/:id/overview",
			security.RequireCoachOf(h.svc, "id"),
			h.GetOverview,
		)
	}
}

type createLinkReq struct {
	DiscipleID string `json:"disciple_id"` // opcional para auto-vínculo
	AutoAccept bool   `json:"auto_accept"` // opcional (para pruebas)
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

func (h *CoachHandler) GetOverview(c *gin.Context) {
	discipleID := c.Param("id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	if days <= 0 {
		days = 14
	}
	metric := c.DefaultQuery("metric", "volume")
	tz := c.DefaultQuery("tz", "America/Santiago")

	coachID := security.MustUserID(c) // <--- en vez de MustClaims
	if coachID == "" {
		return // ya abortó con 401
	}

	out, err := h.svc.GetOverview(c.Request.Context(), coachID, discipleID, days, metric, tz)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
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

func (h *CoachHandler) GetTodayForDisciple(c *gin.Context) {
	discipleID := c.Param("id")
	tz := c.DefaultQuery("tz", "America/Santiago")

	resp, err := h.hist.GetMeTodayFor(c.Request.Context(), discipleID, tz)
	if err != nil {
		// Si tu svc ya convierte ErrNoDay en respuesta vacía, no llegas aquí.
		// Pero por seguridad, devuelve shape vacío en vez de 500:
		if strings.Contains(err.Error(), "no_day") {
			c.JSON(http.StatusOK, gin.H{
				"assignment_id":              nil,
				"day":                        nil,
				"prescriptions":              []interface{}{},
				"current_session_id":         nil,
				"current_session_started_at": nil,
				"current_session_sets_count": nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "detail": err.Error()})
		return
	}

	if resp == nil {
		c.JSON(http.StatusOK, gin.H{
			"assignment_id":              nil,
			"day":                        nil,
			"prescriptions":              []interface{}{},
			"current_session_id":         nil,
			"current_session_started_at": nil,
			"current_session_sets_count": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
