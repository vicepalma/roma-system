package repository

import (
	"context"

	"gorm.io/gorm"
)

type AssignmentDaysItem struct {
	ID                 string   `json:"id"`
	WeekID             string   `json:"week_id"`
	WeekIndex          int      `json:"week_index"`
	DayIndex           int      `json:"day_index"`
	Title              *string  `json:"title,omitempty"`
	Notes              *string  `json:"notes"`
	PrescriptionsCount int      `json:"prescriptions_count"`
	ExerciseNames      []string `json:"exercise_names"`
	IsSessionDay       bool     `json:"is_session_day"`
}

type AssignmentDaysRepository interface {
	LoadAssignmentOwner(ctx context.Context, assignmentID string) (programID string, discipleID string, err error)
	ListAssignmentDays(ctx context.Context, assignmentID, programID string) ([]AssignmentDaysItem, error)
}

type assignmentDaysRepository struct{ db *gorm.DB }

func NewAssignmentDaysRepository(db *gorm.DB) AssignmentDaysRepository {
	return &assignmentDaysRepository{db: db}
}

func (r *assignmentDaysRepository) LoadAssignmentOwner(ctx context.Context, assignmentID string) (string, string, error) {
	const q = `
SELECT a.program_id, a.disciple_id
FROM assignments a
WHERE a.id = $1
LIMIT 1;
`
	var programID, discipleID string
	row := r.db.WithContext(ctx).Raw(q, assignmentID).Row()
	if err := row.Scan(&programID, &discipleID); err != nil {
		return "", "", err
	}
	return programID, discipleID, nil
}

func (r *assignmentDaysRepository) ListAssignmentDays(ctx context.Context, assignmentID, programID string) ([]AssignmentDaysItem, error) {
	// title = COALESCE(d.notes, join de primeras 2-3 exercises)
	const q = `
WITH latest_sess AS (
  SELECT sl.day_id
  FROM session_logs sl
  WHERE sl.assignment_id = $1
  ORDER BY sl.performed_at DESC, sl.created_at DESC
  LIMIT 1
),
days AS (
  SELECT d.id, d.week_id, w.week_index, d.day_index, NULLIF(d.notes,'') AS notes
  FROM program_weeks w
  JOIN program_days d ON d.week_id = w.id
  WHERE w.program_id = $2
)
SELECT
  d.id,
  d.week_id,
  d.week_index,
  d.day_index,
  -- prescriptions_count
  COALESCE((
    SELECT COUNT(*)::int FROM prescriptions p WHERE p.day_id = d.id
  ), 0) AS prescriptions_count,
  -- exercise_names (primeras 3 por orden)
  COALESCE((
    SELECT ARRAY(
      SELECT e.name
      FROM prescriptions p
      JOIN exercises e ON e.id = p.exercise_id
      WHERE p.day_id = d.id
      ORDER BY p.position ASC, p.id ASC
      LIMIT 3
    )
  ), ARRAY[]::text[]) AS exercise_names,
  -- is_session_day
  EXISTS(SELECT 1 FROM latest_sess ls WHERE ls.day_id = d.id) AS is_session_day,
  -- title derivado
  CASE
    WHEN d.notes IS NOT NULL THEN d.notes
    ELSE (
      SELECT CASE
               WHEN COUNT(*) = 0 THEN NULL
               WHEN COUNT(*) = 1 THEN MIN(e.name)
               ELSE array_to_string(ARRAY(
                      SELECT e2.name
                      FROM prescriptions p2
                      JOIN exercises e2 ON e2.id = p2.exercise_id
                      WHERE p2.day_id = d.id
                      ORDER BY p2.position ASC, p2.id ASC
                      LIMIT 2
                    ), ' + ')
             END
      FROM prescriptions p
      JOIN exercises e ON e.id = p.exercise_id
      WHERE p.day_id = d.id
    )
  END AS title
FROM days d
ORDER BY d.week_index ASC, d.day_index ASC;
`
	rows, err := r.db.WithContext(ctx).Raw(q, assignmentID, programID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]AssignmentDaysItem, 0, 32)
	for rows.Next() {
		var it AssignmentDaysItem
		var title *string
		if err := rows.Scan(
			&it.ID,
			&it.WeekID,
			&it.WeekIndex,
			&it.DayIndex,
			&it.PrescriptionsCount,
			&it.ExerciseNames,
			&it.IsSessionDay,
			&title,
		); err != nil {
			return nil, err
		}
		it.Title = title
		// Notes ya viene en title CTE si existía; la queremos también aparte:
		// Volvemos a leer notas con una query chica solo si quieres la nota exacta (opcional).
		out = append(out, it)
	}
	return out, nil
}
