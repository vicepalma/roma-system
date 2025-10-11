package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"gorm.io/gorm"
)

type AuthHandler struct {
	users repository.UserRepository
}

func NewAuthHandler(u repository.UserRepository) *AuthHandler { return &AuthHandler{users: u} }

type registerReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(r *gin.RouterGroup) {
	r.POST("/auth/register", h.register)
	r.POST("/auth/login", h.login)
	r.POST("/auth/refresh", h.refresh)
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash_error"})
		return
	}
	u := &domain.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
	}
	if err := h.users.Create(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email_in_use"})
		return
	}
	tokens, err := security.GenerateTokens(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token_error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": gin.H{"id": u.ID, "name": u.Name, "email": u.Email}, "tokens": tokens})
}

// @Summary Login
// @Description Autentica un usuario y devuelve access/refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body struct{Email string `json:"email"`; Password string `json:"password"`} true "Credenciales"
// @Success 200 {object} struct{Tokens struct{Access string `json:"access"`; Refresh string `json:"refresh"`}; User struct{ID string `json:"id"`; Email string `json:"email"`; Name string `json:"name"`} `json:"user"`}
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]

func (h *AuthHandler) login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	u, err := h.users.FindByEmail(c.Request.Context(), strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	ok, _ := security.CheckPassword(req.Password, u.PasswordHash)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
		return
	}
	tokens, err := security.GenerateTokens(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": u.ID, "name": u.Name, "email": u.Email}, "tokens": tokens})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	type refreshReq struct {
		Refresh string `json:"refresh" binding:"required"`
	}
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	tok, claims, err := security.ParseAndValidate(req.Refresh)
	if err != nil || !tok.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}
	if typ, _ := claims["typ"].(string); typ != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong_token_type"})
		return
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no_sub"})
		return
	}
	tokens, err := security.GenerateTokens(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}
