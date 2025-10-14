package http

import (
	"database/sql"
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
	db    *gorm.DB
}

func NewAuthHandler(u repository.UserRepository, db *gorm.DB) *AuthHandler {
	return &AuthHandler{users: u, db: db}
}

type registerReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // coach|disciple (hint – no se persiste en BD)
}

func (h *AuthHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/auth")
	{
		// Puedes apuntar signup a la misma lógica de register
		g.POST("/signup", h.signup)
		g.POST("/register", h.register)
		g.POST("/login", h.login)
		g.POST("/refresh", h.refresh)
	}
	// /me protegido
	r.GET("/me", security.AuthRequired(), h.me)
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
	c.JSON(http.StatusCreated, gin.H{
		"user":   gin.H{"id": u.ID, "name": u.Name, "email": u.Email},
		"tokens": tokens,
	})
}

func (h *AuthHandler) signup(c *gin.Context) {
	// Reutiliza la misma lógica de register para no duplicar
	var req signupReq
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
	c.JSON(http.StatusOK, gin.H{
		"user":          gin.H{"id": u.ID, "email": u.Email, "name": u.Name},
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

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

func (h *AuthHandler) me(c *gin.Context) {
	uid := security.UserID(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Trae datos básicos del usuario
	u, err := h.users.FindByID(c.Request.Context(), uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}

	// Deriva el rol en base a tu esquema actual (sin columnas nuevas)
	var role string
	err = h.db.WithContext(c.Request.Context()).Raw(`
		SELECT CASE
			WHEN EXISTS (SELECT 1 FROM coach_links cl WHERE cl.coach_id = ? AND cl.status = 'accepted') THEN 'coach'
			WHEN EXISTS (SELECT 1 FROM programs p WHERE p.owner_id = ?) THEN 'coach'
			ELSE 'disciple'
		END AS role;
	`, uid, uid).Row().Scan(&role)
	if err != nil || role == "" {
		role = "disciple"
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"role":  role,
	})
}
