package service

import (
	"time"

	"github.com/888tru/jobtrack-api/internal/repository"
)

type StatsService interface {
	GetStats(userID uint) (*StatsData, error)
}

type StatsData struct {
	Total      int64            `json:"total"`
	ByStatus   map[string]int64 `json:"by_status"`
	ThisMonth  int64            `json:"this_month"`
	Interviews int64            `json:"interviews"`
}

type statsService struct {
	appRepo       repository.JobApplicationRepository
	interviewRepo repository.InterviewRepository
}

func NewStatsService(appRepo repository.JobApplicationRepository, interviewRepo repository.InterviewRepository) StatsService {
	return &statsService{appRepo: appRepo, interviewRepo: interviewRepo}
}

func (s *statsService) GetStats(userID uint) (*StatsData, error) {
	counts, err := s.appRepo.CountByStatus(userID)
	if err != nil {
		return nil, err
	}

	byStatus := make(map[string]int64)
	var total int64
	for _, c := range counts {
		byStatus[string(c.Status)] = c.Count
		total += c.Count
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	thisMonth, err := s.appRepo.CountSince(userID, startOfMonth)
	if err != nil {
		return nil, err
	}

	interviews, err := s.interviewRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &StatsData{
		Total:      total,
		ByStatus:   byStatus,
		ThisMonth:  thisMonth,
		Interviews: interviews,
	}, nil
}
