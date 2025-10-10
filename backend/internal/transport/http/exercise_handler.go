package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/service"
)

type ExerciseHandler struct {
	svc service.ExerciseService
}

func NewExerciseHandler(s service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{svc: s}
}

func (h *ExerciseHandler) Register(r *gin.RouterGroup) {
	r.GET("/exercises", h.list)
}

// GET /api/exercises?q=pecho&tags=hypertrophy,compound&match=any&limit=50&offset=0
func (h *ExerciseHandler) list(c *gin.Context) {
	q := c.Query("q")
	rawTags := c.Query("tags")
	match := c.DefaultQuery("match", "any") // any | all
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var tags []string
	if rawTags != "" {
		tags = strings.Split(rawTags, ",")
	}

	items, total, err := h.svc.Search(c.Request.Context(), q, tags, match, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"q":      q,
		"tags":   tags,
		"match":  match,
	})
}
