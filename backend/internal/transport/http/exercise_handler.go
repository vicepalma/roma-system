package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type ExerciseHandler struct {
	svc service.ExerciseService
}

func NewExerciseHandler(s service.ExerciseService) *ExerciseHandler { return &ExerciseHandler{svc: s} }

func (h *ExerciseHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/exercises")
	{
		g.GET("", h.list)
		g.POST("", h.create)
		g.GET("/:id", h.get)
		g.PUT("/:id", h.update)
		g.DELETE("/:id", h.delete)
	}
}

func (h *ExerciseHandler) list(c *gin.Context) {
	limit := atoiOrZero(c.Query("limit"))
	offset := atoiOrZero(c.Query("offset"))

	f := repository.ExerciseFilter{
		Query:     c.Query("query"),
		Muscle:    c.Query("muscle"),
		Equipment: c.Query("equipment"),
		Limit:     limit,
		Offset:    offset,
	}
	items, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *ExerciseHandler) create(c *gin.Context) {
	var body service.CreateExercise
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	ex, err := h.svc.Create(c.Request.Context(), body)
	if err != nil {
		// manejar unique(lower(name)) -> 409
		if strings.Contains(strings.ToLower(err.Error()), "uq_exercise_name") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "name_already_exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ex)
}

func (h *ExerciseHandler) get(c *gin.Context) {
	id := c.Param("id")
	ex, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, ex)
}

func (h *ExerciseHandler) update(c *gin.Context) {
	id := c.Param("id")
	var body service.UpdateExercise
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "detail": err.Error()})
		return
	}
	ex, err := h.svc.Update(c.Request.Context(), id, body)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "uq_exercise_name") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "name_already_exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ex)
}

func (h *ExerciseHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

func atoiOrZero(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0
	}
	return n
}
