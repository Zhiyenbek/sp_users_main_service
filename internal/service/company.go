package service

import (
	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/Zhiyenbek/sp-users-main-service/internal/repository"
	"go.uber.org/zap"
)

type companyService struct {
	companyRepo repository.CompanyRepository
	cfg         *config.Configs
	logger      *zap.SugaredLogger
}

func NewCompanyService(companyRepo repository.CompanyRepository, cfg *config.Configs, logger *zap.SugaredLogger) *companyService {
	return &companyService{
		companyRepo: companyRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

func (s *companyService) CreateCompany(company *models.Company) (string, error) {
	return s.companyRepo.CreateCompany(company)
}

func (s *companyService) UpdateCompany(company *models.Company) error {
	return s.companyRepo.UpdateCompany(company)
}

func (s *companyService) GetCompany(publicID string) (*models.Company, error) {
	return s.companyRepo.GetCompany(publicID)
}

func (s *companyService) GetCompanies(args *models.SearchArgs) ([]*models.Company, int, error) {
	return s.companyRepo.GetCompanies(args)
}

func (s *companyService) Exists(publicID string) error {
	exists, err := s.companyRepo.Exists(publicID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrPermissionDenied
	}
	return nil
}
