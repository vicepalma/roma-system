package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type CoachService interface {
	CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error)
	UpdateLinkStatus(ctx context.Context, id string, actorID string, action string) (*domain.CoachLink, error)
	ListLinks(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error)
	CanCoach(ctx context.Context, coachID, discipleID string) (bool, error)

	ListDisciples(ctx context.Context, coachID string) ([]repository.DiscipleRow, error)
	AssignProgram(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*repository.AssignmentMinimal, error)

	GetOverview(ctx context.Context, coachID, discipleID string, days int, metric, tz string) (*CoachOverview, error)
}

type coachService struct {
	repo repository.CoachRepository
	hist HistoryService
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

func NewCoachService(repo repository.CoachRepository, hist HistoryService) CoachService {
	return &coachService{repo: repo, hist: hist}
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

	// Reutiliza tu HistoryService (asegúrate de tener estos wrappers)
	me, err := s.hist.GetMeTodayFor(ctx, discipleID)
	if err != nil && !errors.Is(err, ErrNoDay) {
		return nil, err
	}
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
