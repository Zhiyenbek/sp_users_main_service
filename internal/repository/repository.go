package repository

import (
	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	RecruiterRepository
	CandidateRepository
}

type RecruiterRepository interface {
	Exists(publicID string) (bool, error)
	GetRecruiter(publicID string) (*models.Recruiter, error)
}
type CandidateRepository interface {
	GetCandidatesBySearch(searchArgs *models.SearchArgs) ([]*models.Candidate, int, error)
	GetCandidateByPublicID(publicID string) (*models.Candidate, error)
	Exists(publicID string) (bool, error)
	AddSkillsToCandidate(candidateID string, skills []string) error
	UpdateCandidateByID(candidateID string, updateData *models.Candidate) error
	DeleteCandidateByID(candidateID string) error
	DeleteSkillsFromCandidate(candidateID string, skills []string) error
}

func New(db *pgxpool.Pool, cfg *config.Configs, log *zap.SugaredLogger) *Repository {
	return &Repository{
		RecruiterRepository: NewRecruiterRepository(db, cfg.DB, log),
		CandidateRepository: NewCandidateRepository(db, cfg.DB, log),
	}
}
