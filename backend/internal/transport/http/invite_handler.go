package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type createInviteRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name"`
	TTLHours string `json:"ttl_hours"` // opcional, string para flexibilidad
	BaseURL  string `json:"base_url"`  // opcional si quieres sobreescribir el baseURL del service
}

type InviteHandler struct {
	svc service.InviteService
}

func NewInviteHandler(s service.InviteService) *InviteHandler { return &InviteHandler{svc: s} }

func (h *InviteHandler) Register(r *gin.RouterGroup) {
	grp := r.Group("/coach")
	{
		// Crear invitación (requiere autenticación; opcionalmente podrías validar "rol coach")
		grp.POST("/invitations", h.createInvite)

		// Aceptar invitación (requiere autenticación; el "disciple" es el usuario actual)
		grp.POST("/invitations/:code/accept", h.acceptInvite)
	}
}

func (h *InviteHandler) createInvite(c *gin.Context) {
	coachID := security.UserID(c)
	if coachID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req createInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	ttl := 72
	if strings.TrimSpace(req.TTLHours) != "" {
		if v, err := strconv.Atoi(req.TTLHours); err == nil && v > 0 {
			ttl = v
		}
	}

	res, err := h.svc.CreateInvite(c.Request.Context(), coachID, req.Email, req.Name, ttl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Si el service soporta baseURL dinámico, podrías reconstruir aquí con req.BaseURL
	c.JSON(http.StatusCreated, res)
}

func (h *InviteHandler) acceptInvite(c *gin.Context) {
	code := c.Param("code")
	discipleID := security.UserID(c)
	if discipleID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	res, err := h.svc.AcceptInvite(c.Request.Context(), code, discipleID)
	if err != nil {
		switch err.Error() {
		case "invalid_code":
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid_code"})
		case "invite_expired":
			c.JSON(http.StatusGone, gin.H{"error": "invite_expired"})
		case "invite_not_pending":
			c.JSON(http.StatusConflict, gin.H{"error": "invite_not_pending"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, res)
}
