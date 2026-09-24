package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"gorm.io/gorm"
	"time"
)

type TemplateHandler struct{ db *gorm.DB }

func NewTemplateHandler(db *gorm.DB) *TemplateHandler { return &TemplateHandler{db: db} }

type templateItem struct {
	ID          string     `json:"id"`
	OwnerID     string     `json:"owner_id"`
	Title       string     `json:"title"`
	Notes       *string    `json:"notes,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}
type templatePreview struct {
	templateItem
	Weeks []templateWeek `json:"weeks"`
}
type templateWeek struct {
	ID        string        `json:"id"`
	WeekIndex int           `json:"week_index"`
	Days      []templateDay `json:"days"`
}
type templateDay struct {
	ID            string                 `json:"id"`
	DayIndex      int                    `json:"day_index"`
	Title         *string                `json:"title,omitempty"`
	Notes         *string                `json:"notes,omitempty"`
	Prescriptions []templatePrescription `json:"prescriptions"`
}
type templatePrescription struct {
	ID           string   `json:"id"`
	ExerciseID   string   `json:"exercise_id"`
	ExerciseName string   `json:"exercise_name"`
	Series       int      `json:"series"`
	Reps         string   `json:"reps"`
	RestSec      *int     `json:"rest_sec,omitempty"`
	ToFailure    bool     `json:"to_failure"`
	Tempo        *string  `json:"tempo,omitempty"`
	RIR          *int     `json:"rir,omitempty"`
	RPE          *float32 `json:"rpe,omitempty"`
	MethodID     *string  `json:"method_id,omitempty"`
	Notes        *string  `json:"notes,omitempty"`
	Position     int      `json:"position"`
}

func (h *TemplateHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/templates")
	g.GET("", h.list)
	g.POST("", h.create)
	g.GET("/:id", h.get)
	g.POST("/:id/publish", h.publish)
	g.POST("/:id/import", h.importTemplate)
}
func (h *TemplateHandler) role(c *gin.Context) string {
	role, _ := security.RoleOf(h.db.WithContext(c.Request.Context()), userID(c))
	return role
}
func (h *TemplateHandler) list(c *gin.Context) {
	var rows []templateItem
	q := h.db.WithContext(c).Raw(`SELECT id,owner_id,title,notes,status,created_at,published_at FROM program_templates WHERE status='published' OR owner_id=? ORDER BY created_at DESC`, userID(c))
	if err := q.Scan(&rows).Error; err != nil {
		c.JSON(500, gin.H{"error": "server_error"})
		return
	}
	c.JSON(200, gin.H{"items": rows, "total": len(rows)})
}
func (h *TemplateHandler) create(c *gin.Context) {
	if h.role(c) != "coach" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var in struct {
		ProgramID string  `json:"program_id" binding:"required"`
		Title     string  `json:"title" binding:"required,min=2"`
		Notes     *string `json:"notes"`
	}
	if c.ShouldBindJSON(&in) != nil {
		c.JSON(400, gin.H{"error": "bad_request"})
		return
	}
	var owner string
	if err := h.db.WithContext(c).Raw(`SELECT owner_id FROM programs WHERE id=?`, in.ProgramID).Row().Scan(&owner); err != nil || owner != userID(c) {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var out templateItem
	err := h.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`INSERT INTO program_templates(owner_id,title,notes) VALUES(?,?,?) RETURNING id,owner_id,title,notes,status,created_at,published_at`, userID(c), in.Title, in.Notes).Scan(&out).Error; err != nil {
			return err
		}
		var weeks []struct {
			ID        string
			WeekIndex int
		}
		if err := tx.Raw(`SELECT id,week_index FROM program_weeks WHERE program_id=? ORDER BY week_index`, in.ProgramID).Scan(&weeks).Error; err != nil {
			return err
		}
		for _, w := range weeks {
			var nw struct{ ID string }
			if err := tx.Raw(`INSERT INTO program_template_weeks(template_id,week_index) VALUES(?,?) RETURNING id`, out.ID, w.WeekIndex).Scan(&nw).Error; err != nil {
				return err
			}
			var days []struct {
				ID       string
				DayIndex int
				Title    *string
				Notes    *string
			}
			if err := tx.Raw(`SELECT id,day_index,title,notes FROM program_days WHERE week_id=? ORDER BY day_index`, w.ID).Scan(&days).Error; err != nil {
				return err
			}
			for _, d := range days {
				var nd struct{ ID string }
				if err := tx.Raw(`INSERT INTO program_template_days(week_id,day_index,title,notes) VALUES(?,?,?,?) RETURNING id`, nw.ID, d.DayIndex, d.Title, d.Notes).Scan(&nd).Error; err != nil {
					return err
				}
				if err := tx.Exec(`INSERT INTO program_template_prescriptions(day_id,exercise_id,series,reps,rest_sec,to_failure,tempo,rir,rpe,method_id,notes,position) SELECT ?,exercise_id,series,reps,rest_sec,to_failure,tempo,rir,rpe,method_id,notes,position FROM prescriptions WHERE day_id=?`, nd.ID, d.ID).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "server_error"})
		return
	}
	c.JSON(201, out)
}
func (h *TemplateHandler) get(c *gin.Context) {
	var out templatePreview
	if err := h.db.WithContext(c).Raw(`SELECT id,owner_id,title,notes,status,created_at,published_at FROM program_templates WHERE id=? AND (status='published' OR owner_id=?)`, c.Param("id"), userID(c)).Scan(&out.templateItem).Error; err != nil || out.ID == "" {
		c.JSON(404, gin.H{"error": "not_found"})
		return
	}
	var weeks []templateWeek
	h.db.WithContext(c).Raw(`SELECT id,week_index FROM program_template_weeks WHERE template_id=? ORDER BY week_index`, out.ID).Scan(&weeks)
	for i := range weeks {
		h.db.WithContext(c).Raw(`SELECT id,day_index,title,notes FROM program_template_days WHERE week_id=? ORDER BY day_index`, weeks[i].ID).Scan(&weeks[i].Days)
		for j := range weeks[i].Days {
			h.db.WithContext(c).Raw(`SELECT p.id,p.exercise_id,e.name AS exercise_name,p.series,p.reps,p.rest_sec,p.to_failure,p.tempo,p.rir,p.rpe,p.method_id,p.notes,p.position FROM program_template_prescriptions p JOIN exercises e ON e.id=p.exercise_id WHERE p.day_id=? ORDER BY p.position`, weeks[i].Days[j].ID).Scan(&weeks[i].Days[j].Prescriptions)
		}
	}
	out.Weeks = weeks
	c.JSON(200, out)
}
func (h *TemplateHandler) publish(c *gin.Context) {
	if h.role(c) != "coach" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	res := h.db.WithContext(c).Exec(`UPDATE program_templates SET status='published',published_at=now(),updated_at=now() WHERE id=? AND owner_id=?`, c.Param("id"), userID(c))
	if res.Error != nil {
		c.JSON(500, gin.H{"error": "server_error"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "not_found"})
		return
	}
	c.Status(204)
}
func (h *TemplateHandler) importTemplate(c *gin.Context) {
	var role string
	role, _ = security.RoleOf(h.db.WithContext(c.Request.Context()), userID(c))
	if role != "coach" && role != "disciple" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var outID string
	err := h.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var t templateItem
		if err := tx.Raw(`SELECT id,title,notes,status FROM program_templates WHERE id=? AND status='published'`, c.Param("id")).Scan(&t).Error; err != nil || t.ID == "" {
			return gorm.ErrRecordNotFound
		}
		kind := "coach_program"
		if role == "disciple" {
			kind = "self_training"
		}
		var p struct{ ID string }
		if err := tx.Raw(`INSERT INTO programs(owner_id,title,notes,kind) VALUES(?,?,?,?) RETURNING id`, userID(c), t.Title, t.Notes, kind).Scan(&p).Error; err != nil {
			return err
		}
		outID = p.ID
		var weeks []struct {
			ID        string
			WeekIndex int
		}
		tx.Raw(`SELECT id,week_index FROM program_template_weeks WHERE template_id=? ORDER BY week_index`, t.ID).Scan(&weeks)
		for _, w := range weeks {
			var nw struct{ ID string }
			if err := tx.Raw(`INSERT INTO program_weeks(program_id,week_index) VALUES(?,?) RETURNING id`, p.ID, w.WeekIndex).Scan(&nw).Error; err != nil {
				return err
			}
			var days []struct {
				ID       string
				DayIndex int
				Title    *string
				Notes    *string
			}
			tx.Raw(`SELECT id,day_index,title,notes FROM program_template_days WHERE week_id=? ORDER BY day_index`, w.ID).Scan(&days)
			for _, d := range days {
				var nd struct{ ID string }
				if err := tx.Raw(`INSERT INTO program_days(week_id,day_index,title,notes) VALUES(?,?,?,?) RETURNING id`, nw.ID, d.DayIndex, d.Title, d.Notes).Scan(&nd).Error; err != nil {
					return err
				}
				if err := tx.Exec(`INSERT INTO prescriptions(day_id,exercise_id,series,reps,rest_sec,to_failure,tempo,rir,rpe,method_id,notes,position) SELECT ?,exercise_id,series,reps,rest_sec,to_failure,tempo,rir,rpe,method_id,notes,position FROM program_template_prescriptions WHERE day_id=?`, nd.ID, d.ID).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{"error": "not_found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "server_error"})
		return
	}
	kind := "coach_program"
	if role == "disciple" {
		kind = "self_training"
	}
	c.JSON(201, gin.H{"program_id": outID, "kind": kind})
}
