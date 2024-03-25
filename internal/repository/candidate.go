package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zhiyenbek/sp-users-main-service/config"
	"github.com/Zhiyenbek/sp-users-main-service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
			SELECT
				c.public_id,
				c.current_position,
				c.education,
				c.resume,
				c.bio,
				u.photo,
				u.first_name,
				u.last_name,
				array_agg(s.name) AS skills
			FROM
				candidates c
			JOIN
				users u ON c.public_id = u.public_id
			LEFT JOIN
				candidate_skills cs ON c.id = cs.candidate_id
			LEFT JOIN
				skills s ON cs.skill_id = s.id
			WHERE
				LOWER(u.first_name) LIKE LOWER($1) OR LOWER(u.last_name) LIKE LOWER($1)
			GROUP BY
				c.public_id,
				c.current_position,
				c.education,
				c.resume,
				c.bio,
				u.first_name,
				u.last_name,
				u.photo
			OFFSET $2
			LIMIT $3
	`

	countQuery := `
	SELECT COUNT(DISTINCT c.id)
	FROM candidates c
	JOIN users u ON c.public_id = u.public_id
	WHERE LOWER(u.first_name) LIKE LOWER($1) OR LOWER(u.last_name) LIKE LOWER($1)	
	`

	var totalCount int
	fmt.Println(query)
	err := r.db.QueryRow(ctx, countQuery, "%"+searchArgs.Search+"%").Scan(&totalCount)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while fetching candidates count: %v", err)
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, "%"+searchArgs.Search+"%", (searchArgs.PageNum-1)*searchArgs.PageSize, searchArgs.PageSize)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("Error occurred while fetching candidates: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	candidates := make([]*models.Candidate, 0)
	fmt.Println(rows.RawValues())
	for rows.Next() {
		candidate := &models.Candidate{}

		err := rows.Scan(
			&candidate.PublicID,
			&candidate.CurrentPosition,
			&candidate.Education,
			&candidate.Resume,
			&candidate.Bio,
			&candidate.Photo,
			&candidate.FirstName,
			&candidate.LastName,
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
	query := `SELECT c.id, c.public_id, c.current_position, c.education, c.resume, c.bio, u.first_name, u.last_name, u.photo
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
		&result.Photo,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return nil, err
	}

	query = `SELECT array_agg(DISTINCT s.name) from skills s
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

	query = `SELECT array_agg(DISTINCT i.public_id) from interviews i
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

func (r *candidateRepository) UpdateCandidateByID(candidateID string, updateData *models.Candidate) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
	UPDATE candidates
	SET
		current_position = COALESCE($2, current_position),
		education = COALESCE($3, education),
		resume = COALESCE($4, resume),
		bio = COALESCE($5, bio)
	WHERE
		public_id = $1
	`

	_, err := r.db.Exec(ctx, query, candidateID, updateData.CurrentPosition, updateData.Education, updateData.Resume, updateData.Bio)
	if err != nil {
		r.logger.Errorf("Error updating candidate: %v", err)
		return err
	}

	query = `
	UPDATE users
	SET
		first_name = COALESCE($2, first_name),
		last_name = COALESCE($3, last_name),
		photo = COALESCE($4, photo)
	WHERE
		public_id = $1
	`

	_, err = r.db.Exec(ctx, query, candidateID, updateData.FirstName, updateData.LastName, updateData.Photo)
	if err != nil {
		r.logger.Errorf("Error updating candidate's user data: %v", err)
		return err
	}

	return nil
}

func (r *candidateRepository) DeleteCandidateByID(candidateID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
    DELETE FROM candidates
    WHERE public_id = $1
    `

	_, err := r.db.Exec(ctx, query, candidateID)
	if err != nil {
		r.logger.Errorf("Error deleting candidate: %v", err)
		return err
	}

	return nil
}
func (r *candidateRepository) AddSkillsToCandidate(candidateID string, skills []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error beginning transaction: %v", err)
		return err
	}

	// Loop over the skills array
	for _, skillName := range skills {
		// Check if the skill already exists in the database
		var skillID int
		query := `
		SELECT id FROM skills WHERE name = $1
		`
		err := tx.QueryRow(ctx, query, skillName).Scan(&skillID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Skill doesn't exist, so insert it into the database
				insertQuery := `
				INSERT INTO skills (name) VALUES ($1)
				RETURNING id
				`
				err = tx.QueryRow(ctx, insertQuery, skillName).Scan(&skillID)
				if err != nil {
					r.logger.Errorf("Error inserting new skill: %v", err)
					tx.Rollback(ctx)
					return err
				}
			} else {
				r.logger.Errorf("Error checking skill existence: %v", err)
				tx.Rollback(ctx)
				return err
			}
		}

		// Associate the skill with the candidate
		insertQuery := `
		INSERT INTO candidate_skills (candidate_id, skill_id) VALUES (
			(SELECT id FROM candidates WHERE public_id = $1),
			$2
		) ON CONFLICT DO NOTHING
		`
		_, err = tx.Exec(ctx, insertQuery, candidateID, skillID)
		if err != nil {
			r.logger.Errorf("Error adding skill to candidate: %v", err)
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (r *candidateRepository) DeleteSkillsFromCandidate(candidateID string, skills []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error beginning transaction: %v", err)
		return err
	}

	// Rollback the transaction if an error occurs
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Loop over the skills array
	for _, skillName := range skills {
		// Get the skill ID
		var skillID int
		skillQuery := `
		SELECT id FROM skills WHERE name = $1
		`
		err := tx.QueryRow(ctx, skillQuery, skillName).Scan(&skillID)
		if err != nil {
			if err == pgx.ErrNoRows {
				r.logger.Warnf("Skill %s does not exist", skillName)
				continue // Skill doesn't exist, continue to the next skill
			} else {
				r.logger.Errorf("Error retrieving skill ID: %v", err)
				tx.Rollback(ctx)
				return err
			}
		}

		// Delete the skill from the candidate
		deleteQuery := `
		DELETE FROM candidate_skills
		WHERE candidate_id = (SELECT id FROM candidates WHERE public_id = $1)
		AND skill_id = $2
		`
		_, err = tx.Exec(ctx, deleteQuery, candidateID, skillID)
		if err != nil {
			tx.Rollback(ctx)
			r.logger.Errorf("Error deleting skill from candidate: %v", err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (r *candidateRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM candidates WHERE public_id = $1)`

	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *candidateRepository) GetInterviewsByPublicID(publicID string, searchArgs *models.SearchArgs) ([]*models.InterviewResults, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
	SELECT i.public_id, i.results, p.public_id
	FROM interviews i
	INNER JOIN user_interviews ui ON ui.interview_id = i.id
	INNER JOIN candidates c ON c.id = ui.candidate_id
	INNER JOIN positions p ON p.id = ui.position_id
	WHERE c.public_id = $1
	GROUP BY i.public_id, i.results
	LIMIT $2 OFFSET $3;
`
	offset := (searchArgs.PageNum - 1) * searchArgs.PageSize
	rows, err := r.db.Query(ctx, query, publicID, searchArgs.PageSize, offset)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, 0, nil
		}
		r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	res := make([]*models.InterviewResults, 0)
	for rows.Next() {
		var resultBytes []byte
		result := &models.InterviewResults{}
		err = rows.Scan(
			&result.PublicID,
			&resultBytes,
			&result.PositionPublicID,
		)
		if err != nil {
			r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
			return nil, 0, err
		}
		result.RawResult = resultBytes
		res = append(res, result)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating over interview result for position rows: %v", err)
		return nil, 0, err
	}
	query = `
	SELECT COUNT(*)
	FROM interviews i
	INNER JOIN user_interviews ui ON ui.interview_id = i.id
	INNER JOIN candidates c ON c.id = ui.candidate_id
	WHERE c.public_id = $1
`
	var totalCount int
	err = r.db.QueryRow(ctx, query, publicID).Scan(&totalCount)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, 0, nil
		}
		r.logger.Errorf("Error occurred while retrieving position count: %v", err)
		return nil, 0, err
	}
	return res, totalCount, nil

}
