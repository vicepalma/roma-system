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

type SessionDetail struct {
	ID           string     `json:"id"`
	AssignmentID string     `json:"assignment_id"`
	DayID        string     `json:"day_id"`
	StartedAt    time.Time  `json:"started_at"` // = performed_at
	Status       string     `json:"status"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

type SetPatch struct {
	PrescriptionID *string  `json:"prescription_id,omitempty"`
	SetIndex       *int     `json:"set_index,omitempty"`
	Weight         *float64 `json:"weight,omitempty"`
	Reps           *int     `json:"reps,omitempty"`
	RPE            *float64 `json:"rpe,omitempty"`
	ToFailure      *bool    `json:"to_failure,omitempty"`
}

type SessionService interface {
	Start(ctx context.Context, discipleID, assignmentID, dayID string, performedAt *time.Time, notes *string) (*domain.SessionLog, error)
	Get(ctx context.Context, discipleID, sessionID string) (*domain.SessionLog, []domain.SetRow, []repository.CardioSegment, error)
	AddSet(ctx context.Context, discipleID, sessionID, prescriptionID string, setIndex int, weight *float64, reps int, rpe *float32, toFailure bool) (*domain.SetLog, error)
	AddCardio(ctx context.Context, discipleID, sessionID, modality string, minutes int, hrMin, hrMax *int, notes *string) (*repository.CardioSegment, error)
	ListSets(ctx context.Context, actorID, sessionID string, prescriptionID *string, limit, offset int) ([]repositorySetLog, int64, error)

	GetSession(ctx context.Context, id string) (*SessionDetail, error)
	PatchSession(ctx context.Context, id string, performedAt *time.Time, notes *string, status *string, endedAt *time.Time) (*SessionDetail, error)

	UpdateSet(ctx context.Context, setID string, patch SetPatch) error
	DeleteSet(ctx context.Context, setID string) error

	GetActiveOpenSessionForMe(ctx context.Context, discipleID string) (*domain.SessionLog, error)
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

func (s *sessionService) Get(ctx context.Context, discipleID, sessionID string) (*domain.SessionLog, []domain.SetRow, []repository.CardioSegment, error) {
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

func (s *sessionService) UpdateSet(ctx context.Context, setID string, in SetPatch) error {
	patch := map[string]any{}
	if in.PrescriptionID != nil {
		patch["prescription_id"] = *in.PrescriptionID
	}
	if in.SetIndex != nil {
		patch["set_index"] = *in.SetIndex
	}
	if in.Weight != nil {
		patch["weight"] = *in.Weight
	}
	if in.Reps != nil {
		patch["reps"] = *in.Reps
	}
	if in.RPE != nil {
		patch["rpe"] = *in.RPE
	}
	if in.ToFailure != nil {
		patch["to_failure"] = *in.ToFailure
	}
	if len(patch) == 0 {
		return nil
	}
	return s.repo.UpdateSet(ctx, setID, patch)
}

func (s *sessionService) DeleteSet(ctx context.Context, setID string) error {
	return s.repo.DeleteSet(ctx, setID)
}

func (s *sessionService) GetSession(ctx context.Context, id string) (*SessionDetail, error) {
	meta, err := s.repo.GetSessionMeta(ctx, id)
	if err != nil {
		return nil, err
	}
	return &SessionDetail{
		ID:           meta.ID,
		AssignmentID: meta.AssignmentID,
		DayID:        meta.DayID,
		StartedAt:    meta.PerformedAt,
		Status:       meta.Status,
		EndedAt:      meta.EndedAt,
		Notes:        meta.Notes,
	}, nil
}

func (s *sessionService) PatchSession(ctx context.Context, id string, performedAt *time.Time, notes *string, status *string, endedAt *time.Time) (*SessionDetail, error) {
	patch := map[string]any{}

	if performedAt != nil {
		patch["performed_at"] = *performedAt
	}
	if notes != nil {
		patch["notes"] = *notes
	}
	if status != nil {
		v := strings.ToLower(strings.TrimSpace(*status))
		if v != "open" && v != "closed" {
			return nil, errors.New("invalid_status")
		}
		patch["status"] = v

		// reglas simples de consistencia con ended_at
		if v == "closed" && endedAt == nil {
			now := time.Now()
			patch["ended_at"] = now
		}
		if v == "open" {
			patch["ended_at"] = nil
		}
	}
	if endedAt != nil {
		patch["ended_at"] = *endedAt
		// si se setea ended_at explícitamente y no vino status, asumimos closed
		if status == nil {
			patch["status"] = "closed"
		}
	}

	if len(patch) > 0 {
		patch["updated_at"] = time.Now()
		if err := s.repo.UpdateSession(ctx, id, patch); err != nil {
			return nil, err
		}
	}
	return s.GetSession(ctx, id)
}

func (s *sessionService) GetActiveOpenSessionForMe(ctx context.Context, discipleID string) (*domain.SessionLog, error) {
	return s.repo.GetLatestOpenByDisciple(ctx, discipleID)
}
