package service

import (
	"context"
	"errors"

	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type AssignmentDaysService interface {
	List(ctx context.Context, requesterID, assignmentID string) ([]repository.AssignmentDaysItem, error)
}

type assignmentDaysService struct {
	repo     repository.AssignmentDaysRepository
	coachSvc CoachService
}

func NewAssignmentDaysService(r repository.AssignmentDaysRepository, coach CoachService) AssignmentDaysService {
	return &assignmentDaysService{repo: r, coachSvc: coach}
}

func (s *assignmentDaysService) List(ctx context.Context, requesterID, assignmentID string) ([]repository.AssignmentDaysItem, error) {
	progID, discipleID, err := s.repo.LoadAssignmentOwner(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	// auth: coach del discípulo o el propio discípulo
	if requesterID != discipleID {
		ok, err := s.coachSvc.CanCoach(ctx, requesterID, discipleID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("not_found") // 404 por seguridad
		}
	}
	return s.repo.ListAssignmentDays(ctx, assignmentID, progID)
}
