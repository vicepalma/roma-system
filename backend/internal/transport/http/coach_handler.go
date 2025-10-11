package http

import (
	"net/http"
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

type CoachHandler struct{ svc service.CoachService }

func NewCoachHandler(s service.CoachService) *CoachHandler { return &CoachHandler{svc: s} }

func (h *CoachHandler) Register(r *gin.RouterGroup) {
	grp := r.Group("/coach")
	{
		grp.POST("/links", h.createLink)
		grp.PATCH("/links/:id", h.updateLink)
		grp.GET("/links", h.listLinks)

		grp.GET("/disciples", h.listDisciples)
		grp.POST("/assignments", h.assignProgram)
	}
}

type createLinkReq struct {
	DiscipleID string `json:"disciple_id"` // opcional para auto-vínculo
	AutoAccept bool   `json:"auto_accept"` // opcional (para pruebas)
}

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
