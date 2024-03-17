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

func (s *candidatesService) Exists(publicID string) error {
	exists, err := s.candidateRepo.Exists(publicID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrPermissionDenied
	}
	return nil
}

func (s *candidatesService) AddSkillsToCandidate(candidateID string, skills []string) error {
	return s.candidateRepo.AddSkillsToCandidate(candidateID, skills)
}

func (s *candidatesService) UpdateCandidateByID(candidateID string, updateData *models.Candidate) error {
	return s.candidateRepo.UpdateCandidateByID(candidateID, updateData)
}
func (s *candidatesService) DeleteCandidateByID(candidateID string) error {
	return s.candidateRepo.DeleteCandidateByID(candidateID)
}

func (s *candidatesService) DeleteSkillsFromCandidate(candidateID string, skills []string) error {
	return s.candidateRepo.DeleteSkillsFromCandidate(candidateID, skills)
}
