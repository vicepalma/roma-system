package service

import (
	"context"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"gorm.io/gorm"
)

type repositorySetLog = struct {
	ID             string   `json:"id"`
	SessionID      string   `json:"session_id"`
	PrescriptionID string   `json:"prescription_id"`
	SetIndex       int      `json:"set_index"`
	Weight         *float64 `json:"weight,omitempty"`
	Reps           int      `json:"reps"`
	RPE            *float64 `json:"rpe,omitempty"`
	ToFailure      bool     `json:"to_failure"`
	CreatedAt      string   `json:"created_at"`
}

type SessionService interface {
	Start(ctx context.Context, discipleID, assignmentID, dayID string, performedAt *time.Time, notes *string) (*domain.SessionLog, error)
	Get(ctx context.Context, discipleID, sessionID string) (*domain.SessionLog, []domain.SetLog, []repository.CardioSegment, error)
	AddSet(ctx context.Context, discipleID, sessionID, prescriptionID string, setIndex int, weight *float64, reps int, rpe *float32, toFailure bool) (*domain.SetLog, error)
	Update(ctx context.Context, discipleID, sessionID string, performedAt *time.Time, notes *string) error
	AddCardio(ctx context.Context, discipleID, sessionID, modality string, minutes int, hrMin, hrMax *int, notes *string) (*repository.CardioSegment, error)
	ListSets(ctx context.Context, actorID, sessionID string, prescriptionID *string, limit, offset int) ([]repositorySetLog, int64, error)
}

type sessionService struct {
	repo     repository.SessionRepository
	coachSvc CoachService
}

func NewSessionService(repo repository.SessionRepository, coachSvc CoachService) SessionService {
	return &sessionService{repo: repo, coachSvc: coachSvc}
}

func (s *sessionService) Start(ctx context.Context, discipleID, assignmentID, dayID string, performedAt *time.Time, notes *string) (*domain.SessionLog, error) {
	if discipleID == "" || assignmentID == "" || dayID == "" {
		return nil, errors.New("missing required fields")
	}
	sess := &domain.SessionLog{
		AssignmentID: assignmentID,
		DiscipleID:   discipleID,
		DayID:        dayID,
		PerformedAt:  time.Now().UTC(),
		Notes:        notes,
	}
	if performedAt != nil {
		sess.PerformedAt = *performedAt
	}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *sessionService) Get(ctx context.Context, discipleID, sessionID string) (*domain.SessionLog, []domain.SetLog, []repository.CardioSegment, error) {
	sess, err := s.repo.GetSession(ctx, sessionID, discipleID)
	if err != nil {
		return nil, nil, nil, err
	}
	sets, err := s.repo.ListSets(ctx, sessionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, err
	}
	cardio, err := s.repo.ListCardio(ctx, sessionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, err
	}
	return sess, sets, cardio, nil
}

func (s *sessionService) AddSet(ctx context.Context, discipleID, sessionID, prescriptionID string, setIndex int, weight *float64, reps int, rpe *float32, toFailure bool) (*domain.SetLog, error) {
	// verifica pertenencia del session al usuario
	if _, err := s.repo.GetSession(ctx, sessionID, discipleID); err != nil {
		return nil, err
	}
	set := &domain.SetLog{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		SetIndex:       setIndex,
		Weight:         weight,
		Reps:           reps,
		RPE:            rpe,
		ToFailure:      toFailure,
	}
	return set, s.repo.AddSet(ctx, set)
}

func (s *sessionService) Update(ctx context.Context, discipleID, sessionID string, performedAt *time.Time, notes *string) error {
	if _, err := s.repo.GetSession(ctx, sessionID, discipleID); err != nil {
		return err
	}
	return s.repo.UpdateSession(ctx, sessionID, performedAt, notes)
}

func (s *sessionService) AddCardio(ctx context.Context, discipleID, sessionID, modality string, minutes int, hrMin, hrMax *int, notes *string) (*repository.CardioSegment, error) {
	if _, err := s.repo.GetSession(ctx, sessionID, discipleID); err != nil {
		return nil, err
	}
	if minutes <= 0 {
		return nil, errors.New("minutes must be > 0")
	}
	seg := &repository.CardioSegment{
		SessionID:   sessionID,
		Modality:    modality,
		Minutes:     minutes,
		TargetHRMin: hrMin,
		TargetHRMax: hrMax,
		Notes:       notes,
	}
	return seg, s.repo.AddCardio(ctx, seg)
}

func (s *sessionService) ListSets(ctx context.Context, actorID, sessionID string, prescriptionID *string, limit, offset int) ([]repositorySetLog, int64, error) {
	// 1) Cargar sesión y autorizar (misma regla que POST /sessions/:id/sets)
	sess, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, 0, err
	}
	// Autorización ampliada: discípulo dueño o su coach
	if sess.DiscipleID != actorID {
		ok, err := s.coachSvc.CanCoach(ctx, actorID, sess.DiscipleID)
		if err != nil {
			return nil, 0, err
		}
		if !ok {
			return nil, 0, errors.New("forbidden: not allowed to view this session")
		}
	}

	// 2) Listar sets
	items, total, err := s.repo.ListSessionSets(ctx, sessionID, prescriptionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// map a DTO simple (si prefieres devolver domain.SetLog, elimina este bloque)
	out := make([]repositorySetLog, 0, len(items))
	for _, it := range items {
		var wPtr *float64
		if it.Weight != nil {
			v := float64(*it.Weight)
			wPtr = &v
		}
		var rpePtr *float64
		if it.RPE != nil {
			v := float64(*it.RPE)
			rpePtr = &v
		}
		out = append(out, repositorySetLog{
			ID:             it.ID,
			SessionID:      it.SessionID,
			PrescriptionID: it.PrescriptionID,
			SetIndex:       it.SetIndex,
			Weight:         wPtr,
			Reps:           it.Reps,
			RPE:            rpePtr,
			ToFailure:      it.ToFailure,
		})
	}
	return out, total, nil
}
