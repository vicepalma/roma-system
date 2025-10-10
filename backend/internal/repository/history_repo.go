package repository

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type ExerciseCatalogItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DailyExerciseVolume struct {
	Date         string  `json:"date"` // YYYY-MM-DD
	ExerciseID   string  `json:"exercise_id"`
	ExerciseName string  `json:"exercise_name"`
	Volume       float64 `json:"volume"` // sum(reps*weight)
	Sets         int     `json:"sets"`
	Reps         int     `json:"reps"`
}

type DailyMuscleVolume struct {
	Date          string  `json:"date"` // YYYY-MM-DD
	PrimaryMuscle string  `json:"primary_muscle"`
	Volume        float64 `json:"volume"`
	Sets          int     `json:"sets"`
	Reps          int     `json:"reps"`
}

type HistoryRepository interface {
	ListRecentSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SessionLog, error)
	ListSetsInSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SetLog, error)
	BestSetsByExercise(ctx context.Context, discipleID string) ([]PRRow, error)

	DailyVolumeByExercise(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyExerciseVolume, error)
	DailyVolumeByMuscle(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyMuscleVolume, error)

	ListRelevantExercisesForUser(ctx context.Context, discipleID string) ([]ExerciseCatalogItem, error)
}

type historyRepository struct{ db *gorm.DB }

func NewHistoryRepository(db *gorm.DB) HistoryRepository { return &historyRepository{db: db} }

func (r *historyRepository) ListRecentSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SessionLog, error) {
	var rows []domain.SessionLog
	err := r.db.WithContext(ctx).
		Where("disciple_id = ? AND performed_at >= ?", discipleID, since).
		Order("performed_at DESC").
		Find(&rows).Error
	return rows, err
}

func (r *historyRepository) ListSetsInSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SetLog, error) {
	var rows []domain.SetLog
	err := r.db.WithContext(ctx).
		Joins("JOIN session_logs s ON s.id = set_logs.session_id").
		Where("s.disciple_id = ? AND s.performed_at >= ?", discipleID, since).
		Find(&rows).Error
	return rows, err
}

type PRRow struct {
	ExerciseID   string   `json:"exercise_id"`
	MaxWeight    *float64 `json:"max_weight,omitempty"`
	MaxReps      int      `json:"max_reps,omitempty"`
	Estimated1RM *float64 `json:"estimated_1rm,omitempty"`
}

func (r *historyRepository) BestSetsByExercise(ctx context.Context, discipleID string) ([]PRRow, error) {
	rows := []PRRow{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
		  p.exercise_id,
		  MAX(set_logs.weight)::float AS max_weight,
		  MAX(set_logs.reps)         AS max_reps,
		  MAX(
		    CASE 
		      WHEN set_logs.weight IS NOT NULL 
		           AND set_logs.weight > 0
		           AND set_logs.reps BETWEEN 1 AND 36
		      THEN (set_logs.weight::float * (36.0 / NULLIF(37.0 - set_logs.reps, 0)))
		      ELSE NULL
		    END
		  ) AS estimated_1rm
		FROM set_logs
		JOIN session_logs s ON s.id = set_logs.session_id
		JOIN prescriptions p ON p.id = set_logs.prescription_id
		WHERE s.disciple_id = ?
		  AND set_logs.weight IS NOT NULL
		  AND set_logs.reps BETWEEN 1 AND 36
		GROUP BY p.exercise_id
		ORDER BY estimated_1rm DESC NULLS LAST, max_weight DESC NULLS LAST, max_reps DESC
	`, discipleID).Scan(&rows).Error
	return rows, err
}

func (r *historyRepository) DailyVolumeByExercise(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyExerciseVolume, error) {
	rows := []DailyExerciseVolume{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
		  to_char( (s.performed_at AT TIME ZONE ? )::date, 'YYYY-MM-DD') AS date,
		  p.exercise_id,
		  e.name AS exercise_name,
		  SUM(COALESCE(set_logs.reps,0) * COALESCE(set_logs.weight,0))::float AS volume,
		  COUNT(*) AS sets,
		  SUM(COALESCE(set_logs.reps,0)) AS reps
		FROM set_logs
		JOIN session_logs s  ON s.id = set_logs.session_id
		JOIN prescriptions p ON p.id = set_logs.prescription_id
		JOIN exercises e     ON e.id = p.exercise_id
		WHERE s.disciple_id = ?
		  AND (s.performed_at AT TIME ZONE ? )::date >= ?::date
		GROUP BY 1,2,3
		ORDER BY 1 ASC, 3 ASC
	`, tz, discipleID, tz, sinceDate).Scan(&rows).Error
	return rows, err
}

func (r *historyRepository) DailyVolumeByMuscle(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyMuscleVolume, error) {
	rows := []DailyMuscleVolume{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
		  to_char( (s.performed_at AT TIME ZONE ? )::date, 'YYYY-MM-DD') AS date,
		  lower(e.primary_muscle) AS primary_muscle,
		  SUM(COALESCE(set_logs.reps,0) * COALESCE(set_logs.weight,0))::float AS volume,
		  COUNT(*) AS sets,
		  SUM(COALESCE(set_logs.reps,0)) AS reps
		FROM set_logs
		JOIN session_logs s   ON s.id = set_logs.session_id
		JOIN prescriptions p  ON p.id = set_logs.prescription_id
		JOIN exercises e      ON e.id = p.exercise_id
		WHERE s.disciple_id = ?
		  AND (s.performed_at AT TIME ZONE ? )::date >= ?::date
		GROUP BY 1,2
		ORDER BY 1 ASC, 2 ASC
	`, tz, discipleID, tz, sinceDate).Scan(&rows).Error
	return rows, err
}

func (r *historyRepository) ListRelevantExercisesForUser(ctx context.Context, discipleID string) ([]ExerciseCatalogItem, error) {
	rows := []ExerciseCatalogItem{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT e.id, e.name
		FROM assignments a
		JOIN program_weeks  pw ON pw.program_id = a.program_id
		JOIN program_days   pd ON pd.week_id    = pw.id
		JOIN prescriptions  p  ON p.day_id      = pd.id
		JOIN exercises      e  ON e.id          = p.exercise_id
		WHERE a.disciple_id = ?
		  AND a.is_active = true
		ORDER BY e.name ASC
	`, discipleID).Scan(&rows).Error
	return rows, err
}
