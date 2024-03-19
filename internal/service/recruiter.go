package service

import (
	"encoding/json"

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

func (s *recruiterService) GetInterviewsByPublicID(publicID string, searchArgs *models.SearchArgs) ([]*models.InterviewResults, int, error) {
	res, count, err := s.recruiterRepo.GetInterviewsByPublicID(publicID, searchArgs)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if r != nil && r.RawResult != nil {
			result := models.Result{}
			err = json.Unmarshal(r.RawResult, &result)
			if err != nil {
				s.logger.Error(err)
				return nil, 0, err
			}
			r.Result = result
		}
	}
	return res, count, nil

}
