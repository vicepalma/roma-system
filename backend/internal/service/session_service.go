package service

import (
	"context"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"gorm.io/gorm"
)

type SessionService interface {
	Start(ctx context.Context, discipleID, assignmentID, dayID string, performedAt *time.Time, notes *string) (*domain.SessionLog, error)
	Get(ctx context.Context, discipleID, sessionID string) (*domain.SessionLog, []domain.SetLog, []repository.CardioSegment, error)
	AddSet(ctx context.Context, discipleID, sessionID, prescriptionID string, setIndex int, weight *float64, reps int, rpe *float32, toFailure bool) (*domain.SetLog, error)
	Update(ctx context.Context, discipleID, sessionID string, performedAt *time.Time, notes *string) error
	AddCardio(ctx context.Context, discipleID, sessionID, modality string, minutes int, hrMin, hrMax *int, notes *string) (*repository.CardioSegment, error)
}

type sessionService struct{ repo repository.SessionRepository }

func NewSessionService(r repository.SessionRepository) SessionService {
	return &sessionService{repo: r}
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
