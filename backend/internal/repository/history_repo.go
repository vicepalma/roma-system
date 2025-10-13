package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type CurrentSessionInfo struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`
	SetsCount int       `json:"sets_count"`
}

type LatestSessionInfo struct {
	ID        string    `db:"id"`
	StartedAt time.Time `db:"started_at"`
	SetsCount int       `db:"sets_count"`
}

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

type MeTodayDay struct {
	ID       string
	WeekID   string
	DayIndex int
	Notes    sql.NullString
}

type MeTodayPrescription struct {
	ID            string
	DayID         string
	ExerciseID    string
	Series        int
	Reps          string
	RestSec       sql.NullInt32
	ToFailure     bool
	Position      int
	ExerciseName  string
	PrimaryMuscle string
	Equipment     sql.NullString
}

var ErrNoDay = errors.New("no_day")

type HistoryRepository interface {
	ListRecentSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SessionLog, error)
	ListSetsInSessions(ctx context.Context, discipleID string, since time.Time) ([]domain.SetLog, error)
	BestSetsByExercise(ctx context.Context, discipleID string) ([]PRRow, error)

	DailyVolumeByExercise(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyExerciseVolume, error)
	DailyVolumeByMuscle(ctx context.Context, discipleID string, sinceDate string, tz string) ([]DailyMuscleVolume, error)

	ListRelevantExercisesForUser(ctx context.Context, discipleID string) ([]ExerciseCatalogItem, error)

	ResolveToday(ctx context.Context, discipleID string, tz string) (assignmentID string, day *MeTodayDay, prescs []MeTodayPrescription, err error)

	LatestSessionForAssignmentDay(ctx context.Context, assignmentID, dayID string) (*CurrentSessionInfo, error)
	ActiveAssignmentForToday(ctx context.Context, discipleID, tz string) (string, error)
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

func (r *historyRepository) ResolveToday(ctx context.Context, discipleID string, tz string) (string, *MeTodayDay, []MeTodayPrescription, error) {
	// 1) assignment activo más reciente a fecha de HOY en TZ
	const qAssign = `
SELECT a.id, a.program_id
FROM assignments a
WHERE a.disciple_id = $1
  AND a.is_active = true
  AND a.start_date <= (CURRENT_DATE AT TIME ZONE $2)
  AND (a.end_date IS NULL OR a.end_date >= (CURRENT_DATE AT TIME ZONE $2))
ORDER BY a.created_at DESC
LIMIT 1;
`
	var assignID, programID string
	if err := r.db.WithContext(ctx).Raw(qAssign, discipleID, tz).Row().Scan(&assignID, &programID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, nil, ErrNoDay
		}
		return "", nil, nil, err
	}

	// 2) (MVP) Week 1 / Day 1 del programa
	const qDay = `
SELECT d.id, d.week_id, d.day_index, d.notes
FROM program_days d
JOIN program_weeks w ON w.id = d.week_id
WHERE w.program_id = $1 AND w.week_index = 1 AND d.day_index = 1
ORDER BY d.id ASC
LIMIT 1;
`
	var day MeTodayDay
	if err := r.db.WithContext(ctx).Raw(qDay, programID).Row().Scan(&day.ID, &day.WeekID, &day.DayIndex, &day.Notes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, nil, ErrNoDay
		}
		return "", nil, nil, err
	}

	// 3) Prescripciones + datos de ejercicio
	const qPresc = `
SELECT p.id, p.day_id, p.exercise_id, p.series, p.reps, p.rest_sec, p.to_failure, p.position,
       e.name, e.primary_muscle, e.equipment
FROM prescriptions p
JOIN exercises e ON e.id = p.exercise_id
WHERE p.day_id = $1
ORDER BY p.position ASC, p.id ASC;
`
	rows, err := r.db.Raw(qPresc, day.ID).Rows()
	if err != nil {
		return assignID, &day, nil, err
	}
	defer rows.Close()

	out := make([]MeTodayPrescription, 0, 8)
	for rows.Next() {
		var pr MeTodayPrescription
		if err := rows.Scan(
			&pr.ID, &pr.DayID, &pr.ExerciseID, &pr.Series, &pr.Reps, &pr.RestSec, &pr.ToFailure, &pr.Position,
			&pr.ExerciseName, &pr.PrimaryMuscle, &pr.Equipment,
		); err != nil {
			return assignID, &day, nil, err
		}
		out = append(out, pr)
	}
	if err := rows.Err(); err != nil {
		return assignID, &day, nil, err
	}

	return assignID, &day, out, nil
}

func (r *historyRepository) LatestSessionForAssignmentDay(ctx context.Context, assignmentID, dayID string) (*CurrentSessionInfo, error) {
	const q = `
		SELECT s.id,
		       s.performed_at        AS started_at,
		       COALESCE(cnt.c, 0)    AS sets_count
		FROM session_logs s
		LEFT JOIN (
		  SELECT session_id, COUNT(*) AS c
		  FROM set_logs
		  GROUP BY session_id
		) cnt ON cnt.session_id = s.id
		WHERE s.assignment_id = ? AND s.day_id = ?
		ORDER BY s.created_at DESC
		LIMIT 1;
	`
	var res CurrentSessionInfo
	tx := r.db.WithContext(ctx).Raw(q, assignmentID, dayID).Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil // no hay sesión vigente para ese (assignment_id, day_id)
	}
	return &res, nil
}

func (r *historyRepository) ActiveAssignmentForToday(ctx context.Context, discipleID, tz string) (string, error) {
	const q = `
SELECT a.id
FROM assignments a
WHERE a.disciple_id = ?
  AND a.is_active = true
  AND a.start_date <= (CURRENT_DATE AT TIME ZONE ?)
  AND (a.end_date IS NULL OR a.end_date >= (CURRENT_DATE AT TIME ZONE ?))
ORDER BY a.created_at DESC
LIMIT 1;`
	var id sql.NullString
	if err := r.db.WithContext(ctx).Raw(q, discipleID, tz, tz).Scan(&id).Error; err != nil {
		return "", err
	}
	if !id.Valid {
		return "", sql.ErrNoRows
	}
	return id.String, nil
}
