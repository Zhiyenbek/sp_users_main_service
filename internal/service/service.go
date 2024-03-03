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
}

type Service struct {
	CandidatesService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		CandidatesService: NewCandidatesService(repos, cfg, log),
	}
}
