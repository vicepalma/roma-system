package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"gorm.io/gorm"
)

type CalendarDay struct {
	Date  time.Time `json:"date"`
	DayID string    `json:"day_id"`
	Index int       `json:"day_index"`
	Notes *string   `json:"notes,omitempty"`
}

var ErrAssignmentNotFound = errors.New("assignment_not_for_disciple")

type CoachService interface {
	CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error)
	UpdateLinkStatus(ctx context.Context, id string, actorID string, action string) (*domain.CoachLink, error)
	ListLinks(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error)
	CanCoach(ctx context.Context, coachID, discipleID string) (bool, error)

	ListDisciples(ctx context.Context, coachID string) ([]repository.DiscipleRow, error)
	AssignProgram(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*repository.AssignmentMinimal, error)

	GetOverview(ctx context.Context, coachID, discipleID string, days int, metric, tz string) (*CoachOverview, error)

	ListAssignments(ctx context.Context, coachID string, discipleID *string, limit, offset int) ([]repository.AssignmentListRow, int64, error)

	UpdateAssignment(ctx context.Context, id string, endDate *time.Time, isActive *bool) (*repository.AssignmentRow, error)
	AssignmentCalendar(ctx context.Context, id string, from, to time.Time) ([]CalendarDay, error)
	ActivateAssignment(ctx context.Context, discipleID, assignmentID string) error
	GetActiveAssignment(ctx context.Context, discipleID string) (*domain.Assignment, error)
}

type coachService struct {
	db     *gorm.DB
	repo   repository.CoachRepository
	hist   HistoryService
	assign repository.AssignmentRepository
}

type CoachOverview struct {
	DiscipleID string             `json:"disciple_id"`
	MeToday    *MeTodayResponse   `json:"me_today"`
	Pivot      *PivotResponse     `json:"pivot"`
	Adherence  *AdherenceResponse `json:"adherence"`
}

type AdherenceResponse struct {
	DaysRequested int     `json:"days"`
	DaysWithSets  int     `json:"days_with_sets"`
	Rate          float64 `json:"rate"`
}

func NewCoachService(r repository.CoachRepository, hist HistoryService, opts ...any) CoachService {
	var db *gorm.DB
	var ar repository.AssignmentRepository
	for _, o := range opts {
		if v, ok := o.(*gorm.DB); ok {
			db = v
		}
		if v, ok := o.(repository.AssignmentRepository); ok {
			ar = v
		}
	}
	return &coachService{db: db, repo: r, hist: hist, assign: ar}
}

func (s *coachService) CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error) {
	if coachID == "" || discipleID == "" {
		return nil, errors.New("coach_id and disciple_id required")
	}

	// evita duplicados conocidos del lado del coach
	incoming, _, err := s.repo.ListLinksForUser(ctx, coachID)
	if err == nil {
		for _, l := range incoming {
			if l.CoachID == discipleID || l.DiscipleID == discipleID {
				return &l, nil
			}
		}
	}

	return s.repo.CreateLink(ctx, coachID, discipleID, autoAccept)
}

func (s *coachService) UpdateLinkStatus(ctx context.Context, id string, actorID string, action string) (*domain.CoachLink, error) {
	action = strings.ToLower(action)
	if action != "accept" && action != "reject" {
		return nil, errors.New("invalid action")
	}

	// Solo el DISCÍPULO puede aceptar/rechazar: buscamos invitaciones donde actor es el discípulo (incoming)
	incoming, _, err := s.repo.ListLinksForUser(ctx, actorID)
	if err != nil {
		return nil, err
	}

	var target *domain.CoachLink
	for _, l := range incoming { // incoming = soy el DISCÍPULO en estos links
		if l.ID == id {
			target = &l
			break
		}
	}
	if target == nil {
		return nil, errors.New("forbidden: only disciple can update link")
	}

	newStatus := "rejected"
	if action == "accept" {
		newStatus = "accepted"
	}
	return s.repo.UpdateStatus(ctx, id, newStatus, actorID)
}

func (s *coachService) ListLinks(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error) {
	return s.repo.ListLinksForUser(ctx, userID)
}

func (s *coachService) CanCoach(ctx context.Context, coachID, discipleID string) (bool, error) {
	return s.repo.CanCoach(ctx, coachID, discipleID)
}

func (s *coachService) ListDisciples(ctx context.Context, coachID string) ([]repository.DiscipleRow, error) {
	return s.repo.ListDisciples(ctx, coachID)
}

func (s *coachService) AssignProgram(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*repository.AssignmentMinimal, error) {
	// Autorización: el coach debe estar vinculado con el discípulo (o ser él mismo)
	ok, err := s.repo.CanCoach(ctx, coachID, discipleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("forbidden: not a coach of this disciple")
	}

	if programID == "" || discipleID == "" {
		return nil, errors.New("program_id and disciple_id required")
	}
	if startDate.IsZero() {
		startDate = time.Now()
	}
	return s.repo.CreateAssignment(ctx, coachID, discipleID, programID, startDate)
}

func (s *coachService) GetOverview(ctx context.Context, coachID, discipleID string, days int, metric, tz string) (*CoachOverview, error) {
	ok, err := s.repo.CanCoach(ctx, coachID, discipleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("forbidden: not a coach of disciple")
	}

	me, err := s.hist.GetMeTodayFor(ctx, discipleID, tz)
	if err != nil && !errors.Is(err, ErrNoDay) {
		return nil, err
	}
	// pivot y adherence deben poder devolver vacío sin error
	pivot, err := s.hist.GetPivotByExerciseFor(ctx, discipleID, days, metric, tz, true)
	if err != nil {
		return nil, err
	}
	ad, err := s.hist.GetAdherence(ctx, discipleID, days, tz)
	if err != nil {
		return nil, err
	}

	return &CoachOverview{
		DiscipleID: discipleID,
		MeToday:    me,
		Pivot:      pivot,
		Adherence: &AdherenceResponse{
			DaysRequested: days,
			DaysWithSets:  ad.DaysWithSets,
			Rate:          float64(ad.DaysWithSets) / float64(max(1, days)),
		},
	}, nil
}

func (s *coachService) ListAssignments(ctx context.Context, coachID string, discipleID *string, limit, offset int) ([]repository.AssignmentListRow, int64, error) {
	// Si piden filtrar por disciple_id, valida autorización explícita:
	if discipleID != nil && *discipleID != "" {
		ok, err := s.repo.CanCoach(ctx, coachID, *discipleID)
		if err != nil {
			return nil, 0, err
		}
		if !ok {
			return nil, 0, errors.New("forbidden: not a coach of this disciple")
		}
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListAssignmentsForCoach(ctx, coachID, discipleID, limit, offset)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *coachService) UpdateAssignment(ctx context.Context, id string, endDate *time.Time, isActive *bool) (*repository.AssignmentRow, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	patch := map[string]any{}
	if endDate != nil {
		patch["end_date"] = *endDate
	}
	if isActive != nil {
		patch["is_active"] = *isActive
	}
	if len(patch) == 0 {
		return nil, errors.New("nothing to update")
	}
	if err := s.repo.UpdateAssignment(ctx, id, patch); err != nil {
		return nil, err
	}
	return s.repo.GetAssignmentByID(ctx, id)
}

// MVP: asigna los días del programa (week 1) en ciclo sobre el rango [from..to]
func (s *coachService) AssignmentCalendar(ctx context.Context, id string, from, to time.Time) ([]CalendarDay, error) {
	if to.Before(from) {
		from, to = to, from
	}
	asg, err := s.repo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Limitar por start/end del assignment
	start := asg.StartDate
	if start.After(from) {
		from = start
	}
	if asg.EndDate != nil && asg.EndDate.Before(to) {
		to = *asg.EndDate
	}
	if to.Before(from) {
		return []CalendarDay{}, nil
	}

	days, err := s.repo.ListProgramDaysByProgramWeek(ctx, asg.ProgramID, 1)
	if err != nil {
		return nil, err
	}
	if len(days) == 0 {
		return []CalendarDay{}, nil
	}

	// Construimos calendario cíclico
	out := make([]CalendarDay, 0, 32)
	// punto de arranque: desplazamiento desde start al "from"
	diff := int(from.Sub(start).Hours() / 24) // días
	idx := diff % len(days)
	if idx < 0 {
		idx += len(days)
	}
	cur := from

	for !cur.After(to) {
		d := days[idx]
		out = append(out, CalendarDay{
			Date:  cur,
			DayID: d.ID,
			Index: d.DayIndex,
			Notes: d.Notes,
		})
		// siguiente día
		cur = cur.AddDate(0, 0, 1)
		idx = (idx + 1) % len(days)
	}
	return out, nil
}

func (s *coachService) ActivateAssignment(ctx context.Context, discipleID, assignmentID string) error {
	// Verificar pertenencia primero
	var count int64
	if err := s.db.WithContext(ctx).
		Table("assignments").
		Where("id = ? AND disciple_id = ?", assignmentID, discipleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrAssignmentNotFound
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Desactivar otros activos del discípulo
		if err := tx.Table("assignments").
			Where("disciple_id = ? AND is_active = TRUE AND id <> ?", discipleID, assignmentID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Activar éste (y garantizar start_date si es NULL, limpiar end_date)
		if err := tx.Table("assignments").
			Where("id = ? AND disciple_id = ?", assignmentID, discipleID).
			Updates(map[string]any{
				"is_active":  true,
				"end_date":   nil,
				"start_date": gorm.Expr("COALESCE(start_date, CURRENT_DATE)"),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *coachService) GetActiveAssignment(ctx context.Context, discipleID string) (*domain.Assignment, error) {
	return s.repo.GetActiveAssignment(ctx, discipleID)
}
