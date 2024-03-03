package service

import (
	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/Zhiyenbek/sp_users_main_service/internal/repository"
	"go.uber.org/zap"
)

type candidatesService struct {
	cfg           *config.Configs
	logger        *zap.SugaredLogger
	candidateRepo repository.CandidateRepository
}

func NewCandidatesService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *candidatesService {
	return &candidatesService{
		candidateRepo: repo.CandidateRepository,
		cfg:           cfg,
		logger:        logger,
	}
}
func (s *candidatesService) GetCandidatesBySearch(req *models.SearchArgs) ([]*models.Candidate, int, error) {
	res, count, err := s.candidateRepo.GetCandidatesBySearch(req)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}
func (s *candidatesService) GetCandidateByPublicID(publicID string) (*models.Candidate, error) {
	return s.candidateRepo.GetCandidateByPublicID(publicID)
}
