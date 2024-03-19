package service

import (
	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/Zhiyenbek/sp-users-main-service/internal/repository"
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
	GetInterviewsByPublicID(publicID string, searchArgs *models.SearchArgs) ([]*models.InterviewResults, int, error)
}
type RecruiterService interface {
	Exists(publicID string) error
	GetRecruiter(publicID string) (*models.Recruiter, error)
	GetInterviewsByPublicID(publicID string, searchArgs *models.SearchArgs) ([]*models.InterviewResults, int, error)
}

type CompanyService interface {
	CreateCompany(company *models.Company) error
	UpdateCompany(company *models.Company) error
	GetCompany(publicID string) (*models.Company, error)
	GetCompanies(args *models.SearchArgs) ([]*models.Company, int, error)
	Exists(publicID string) error
}
type Service struct {
	CandidatesService
	RecruiterService
	CompanyService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		CandidatesService: NewCandidatesService(repos, cfg, log),
		RecruiterService:  NewRecruitersService(repos, cfg, log),
		CompanyService:    NewCompanyService(repos.CompanyRepository, cfg, log),
	}
}
