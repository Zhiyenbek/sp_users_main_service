package service

import (
	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/Zhiyenbek/sp_users_main_service/internal/repository"
	"go.uber.org/zap"
)

type CandidatesService interface {
	GetCandidatesBySearch(*models.SearchArgs) ([]*models.Candidate, int, error)
	GetCandidateByPublicID(publicID string) (*models.Candidate, error)
	Exists(publicID string) error
	AddSkillsToCandidate(candidateID string, skills []string) error
	UpdateCandidateByID(candidateID string, updateData *models.Candidate) error
	DeleteCandidateByID(candidateID string) error
	DeleteSkillsFromCandidate(candidateID string, skills []string) error
}
type RecruiterService interface {
	Exists(publicID string) error
	GetRecruiter(publicID string) (*models.Recruiter, error)
}
type Service struct {
	CandidatesService
	RecruiterService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		CandidatesService: NewCandidatesService(repos, cfg, log),
		RecruiterService:  NewRecruitersService(repos, cfg, log),
	}
}
