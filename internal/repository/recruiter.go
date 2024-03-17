package repository

import (
	"context"

	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
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

func (r *recruiterRepository) GetRecruiter(publicID string) (*models.Recruiter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	// Retrieve the recruiter's information
	recruiterQuery := `SELECT r.public_id, r.company_public_id, u.first_name, u.last_name, u.photo
	FROM recruiters r
	JOIN users u ON r.public_id = u.public_id
	WHERE r.public_id = $1`

	recruiter := &models.Recruiter{}
	err := r.db.QueryRow(ctx, recruiterQuery, publicID).Scan(
		&recruiter.PublicID,
		&recruiter.CompanyPublicID,
		&recruiter.FirstName,
		&recruiter.LastName,
		&recruiter.Photo,
	)

	if err != nil {
		r.logger.Errorf("Error occurred while retrieving recruiter information: %v", err)
		return nil, err
	}

	// Retrieve the company information
	companyQuery := `SELECT c.public_id, c.name, c.description
	FROM companies c
	WHERE c.public_id = $1`

	company := &models.Company{}
	err = r.db.QueryRow(ctx, companyQuery, recruiter.CompanyPublicID).Scan(
		&company.PublicID,
		&company.Name,
		&company.Description,
	)

	if err != nil {
		r.logger.Errorf("Error occurred while retrieving company information: %v", err)
		return nil, err
	}

	// Retrieve all positions for the recruiter
	positionsQuery := `SELECT p.public_id, p.name, p.status
	FROM positions p
	WHERE p.recruiters_public_id = $1`

	rows, err := r.db.Query(ctx, positionsQuery, recruiter.PublicID)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving positions for the recruiter: %v", err)
		return nil, err
	}
	defer rows.Close()

	positions := make([]models.Position, 0)
	for rows.Next() {
		position := models.Position{}
		err := rows.Scan(
			&position.PublicID,
			&position.Name,
			&position.Status,
		)
		if err != nil {
			r.logger.Errorf("Error occurred while scanning position rows: %v", err)
			return nil, err
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating over position rows: %v", err)
		return nil, err
	}

	recruiter.Company = company
	recruiter.Positions = positions

	return recruiter, nil
}

func (r *recruiterRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `SELECT EXISTS (SELECT 1 FROM recruiters WHERE public_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking recruiter existence: %v", err)
		return false, err
	}

	return exists, nil
}
