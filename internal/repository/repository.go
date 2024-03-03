package repository

import (
	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	RecruiterRepository
	CandidateRepository
}

type RecruiterRepository interface {
}
type CandidateRepository interface {
	GetCandidatesBySearch(searchArgs *models.SearchArgs) ([]*models.Candidate, int, error)
	GetCandidateByPublicID(publicID string) (*models.Candidate, error)
}

func New(db *pgxpool.Pool, cfg *config.Configs, log *zap.SugaredLogger) *Repository {
	return &Repository{
		RecruiterRepository: NewRecruiterRepository(db, cfg.DB, log),
		CandidateRepository: NewCandidateRepository(db, cfg.DB, log),
	}
}