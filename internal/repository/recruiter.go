package repository

import (
	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type recruiterRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewRecruiterRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) RecruiterRepository {
	return &recruiterRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
