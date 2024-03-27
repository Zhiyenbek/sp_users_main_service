package repository

import (
	"context"
	"database/sql"

	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type companyRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewCompanyRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) CompanyRepository {
	return &companyRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

// CreateCompany creates a new company in the database
func (r *companyRepository) CreateCompany(company *models.Company) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	var publicID string
	query := `
		INSERT INTO companies (name, logo, description)
		VALUES ($1, $2, $3) RETURNING public_id`

	err := r.db.QueryRow(ctx, query, company.Name, company.Logo, company.Description).Scan(&publicID)
	if err != nil {
		r.logger.Errorf("Error occurred while creating company: %v", err)
		return "", err
	}

	return publicID, nil
}

// UpdateCompany updates an existing company in the database
func (r *companyRepository) UpdateCompany(company *models.Company) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		UPDATE companies
		SET name = COALESCE($2, name), logo = COALESCE($3, logo), description = COALESCE($4, description)
		WHERE public_id = $1`

	_, err := r.db.Exec(ctx, query, company.PublicID, company.Name, company.Logo, company.Description)
	if err != nil {

		r.logger.Errorf("Error occurred while updating company: %v", err)
		return err
	}

	return nil
}

// GetCompany retrieves a company from the database by its public ID
func (r *companyRepository) GetCompany(publicID string) (*models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		SELECT id, public_id, name, logo, description
		FROM companies
		WHERE public_id = $1`

	row := r.db.QueryRow(ctx, query, publicID)

	company := &models.Company{}
	err := row.Scan(&company.ID, &company.PublicID, &company.Name, &company.Logo, &company.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Company not found
		}
		r.logger.Errorf("Error occurred while retrieving company: %v", err)
		return nil, err
	}

	return company, nil
}

// GetCompanies retrieves a list of companies from the database based on search parameters
// along with the total count of companies that match the search criteria
func (r *companyRepository) GetCompanies(args *models.SearchArgs) ([]*models.Company, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		SELECT id, public_id, name, logo, description
		FROM companies
		WHERE name ILIKE $1
		ORDER BY id
		LIMIT $2 OFFSET $3`

	countQuery := `
		SELECT COUNT(*)
		FROM companies
		WHERE name ILIKE $1`

	searchPattern := "%" + args.Search + "%"
	offset := (args.PageNum - 1) * args.PageSize

	rows, err := r.db.Query(ctx, query, searchPattern, args.PageSize, offset)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving companies: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	companies := []*models.Company{}
	for rows.Next() {
		company := &models.Company{}
		err := rows.Scan(&company.ID, &company.PublicID, &company.Name, &company.Logo, &company.Description)
		if err != nil {
			r.logger.Errorf("Error occurred while scanning company: %v", err)
			return nil, 0, err
		}
		companies = append(companies, company)
	}

	// Get the total count of companies
	var totalCount int
	err = r.db.QueryRow(ctx, countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving total count of companies: %v", err)
		return nil, 0, err
	}

	return companies, totalCount, nil
}

// Exists checks if a company with the given public ID exists in the database
func (r *companyRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM companies
			WHERE public_id = $1
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking company existence: %v", err)
		return false, err
	}

	return exists, nil
}
