package service

import (
	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/Zhiyenbek/sp-users-main-service/internal/repository"
	"go.uber.org/zap"
)

type recruiterService struct {
	cfg           *config.Configs
	logger        *zap.SugaredLogger
	recruiterRepo repository.RecruiterRepository
}

func NewRecruitersService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *recruiterService {
	return &recruiterService{
		recruiterRepo: repo.RecruiterRepository,
		cfg:           cfg,
		logger:        logger,
	}
}
func (r *recruiterService) Exists(publicID string) error {
	exists, err := r.recruiterRepo.Exists(publicID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrPermissionDenied
	}
	return nil
}

func (r *recruiterService) GetRecruiter(publicID string) (*models.Recruiter, error) {
	return r.recruiterRepo.GetRecruiter(publicID)
}
