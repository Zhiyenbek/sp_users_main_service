package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zhiyenbek/sp_users_main_service/config"
	"github.com/Zhiyenbek/sp_users_main_service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// candidateRepository represents the repository for managing candidates in the database.
type candidateRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

// NewCandidateRepository creates a new instance of candidateRepository.
func NewCandidateRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) CandidateRepository {
	return &candidateRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *candidateRepository) GetCandidatesBySearch(searchArgs *models.SearchArgs) ([]*models.Candidate, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	query := `
		SELECT c.public_id, c.current_position, c.resume, c.bio, c.education, array_agg(s.name)
		FROM candidates AS c
		INNER JOIN users AS u ON c.public_id = u.public_id
		LEFT JOIN candidate_skills AS cs ON c.id = cs.candidate_id
		LEFT JOIN skills AS s ON cs.skill_id = s.id
		WHERE (LOWER(u.first_name) LIKE LOWER($1) OR LOWER(u.last_name) LIKE LOWER($1))
	`

	countQuery := `
		SELECT COUNT(DISTINCT c.public_id)
		FROM candidates AS c
		INNER JOIN users AS u ON c.public_id = u.public_id
		LEFT JOIN candidate_skills AS cs ON c.id = cs.candidate_id
		LEFT JOIN skills AS s ON cs.skill_id = s.id
		WHERE (LOWER(u.first_name) LIKE LOWER($1) OR LOWER(u.last_name) LIKE LOWER($1))
	`

	params := []interface{}{"%" + searchArgs.Search + "%"}

	// Add skills filter if provided
	if len(searchArgs.Skills) > 0 {
		query += `
		AND c.id IN (
			SELECT candidate_id
			FROM candidate_skills
			INNER JOIN skills ON candidate_skills.skill_id = skills.id
			WHERE skills.name = ANY($2)
			GROUP BY candidate_id
			HAVING COUNT(DISTINCT skills.name) = $3
		)
	`
		countQuery += `
		AND c.id IN (
			SELECT candidate_id
			FROM candidate_skills
			INNER JOIN skills ON candidate_skills.skill_id = skills.id
			WHERE skills.name = ANY($2)
			GROUP BY candidate_id
			HAVING COUNT(DISTINCT skills.name) = $3
		)
	`
		params = append(params, pq.Array(searchArgs.Skills), len(searchArgs.Skills))
	}

	query += `
		GROUP BY c.public_id, c.current_position, c.resume, c.bio, c.education
		ORDER BY c.public_id
		OFFSET $2
		LIMIT $3
	`

	var totalCount int
	fmt.Println(countQuery)
	err := r.db.QueryRow(ctx, countQuery, params...).Scan(&totalCount)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while fetching candidates count: %v", err)
		return nil, 0, err
	}
	fmt.Println(searchArgs.Skills)
	params = append(params, (searchArgs.PageNum-1)*searchArgs.PageSize, searchArgs.PageSize)

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while fetching candidates: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	candidates := make([]*models.Candidate, 0)

	for rows.Next() {
		candidate := &models.Candidate{}

		err := rows.Scan(
			&candidate.PublicID,
			&candidate.CurrentPosition,
			&candidate.Resume,
			&candidate.Bio,
			&candidate.Education,
			&candidate.Skills,
		)

		if err != nil {
			r.logger.Errorf("Error occurred while scanning candidate: %v", err)
			return nil, 0, err
		}

		candidates = append(candidates, candidate)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating through candidate rows: %v", err)
		return nil, 0, err
	}

	return candidates, totalCount, nil
}

func (r *candidateRepository) GetCandidateByPublicID(publicID string) (*models.Candidate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	var candidateID int
	result := &models.Candidate{}
	query := `SELECT c.id, c.public_id, c.current_position, c.education, c.resume, c.bio, u.first_name, u.last_name
	FROM candidates c
	JOIN users u ON c.public_id = u.public_id
	WHERE c.public_id = $1`

	err := r.db.QueryRow(ctx, query, publicID).Scan(
		&candidateID,
		&result.PublicID,
		&result.CurrentPosition,
		&result.Education,
		&result.Resume,
		&result.Bio,
		&result.FirstName,
		&result.LastName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return nil, err
	}

	query = `SELECT array_agg(s.name) from skills s
	INNER JOIN candidate_skills cs ON cs.skill_id = s.id
	INNER JOIN candidates c ON cs.candidate_id = c.id
	WHERE cs.candidate_id = $1`
	err = r.db.QueryRow(ctx, query, candidateID).Scan(
		&result.Skills,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return nil, err
	}

	query = `SELECT array_agg(i.public_id) from interviews i
	INNER JOIN user_interviews ui ON ui.interview_id = i.id
	INNER JOIN candidates c ON ui.candidate_id = c.id
	WHERE ui.candidate_id = $1`
	interviewIDs := make([]string, 0)
	err = r.db.QueryRow(ctx, query, candidateID).Scan(
		&interviewIDs,
	)
	for i := range interviewIDs {
		result.Interviews = append(result.Interviews, models.Interview{
			PublicID: interviewIDs[i],
		})
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return nil, err
	}

	return result, nil
}
