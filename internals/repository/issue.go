package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anshul1310/tf-register/internals/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IssueRepository struct {
	databasePool *pgxpool.Pool
}

func NewIssueRepository(databasePool *pgxpool.Pool) *IssueRepository {
	return &IssueRepository{databasePool: databasePool}
}

func (r *IssueRepository) CreateIssue(ctx context.Context, issue *models.Issue) (*models.Issue, error) {
	if issue == nil {
		return nil, errors.New("issue cannot be nil")
	}
	if issue.IssueID == uuid.Nil {
		issue.IssueID = uuid.New()
	}

	query := `
		INSERT INTO issues (
			issue_id, team_id, user_id, user_name, user_email,
			issue_type, title, description, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, 'open'), now(), now())
		RETURNING issue_id, team_id, user_id, user_name, user_email,
		          issue_type, title, description, status, created_at, updated_at;
	`

	created := &models.Issue{}
	err := r.databasePool.QueryRow(
		ctx,
		query,
		issue.IssueID,
		issue.TeamID,
		issue.UserID,
		issue.UserName,
		issue.UserEmail,
		issue.IssueType,
		issue.Title,
		issue.Description,
		issue.Status,
	).Scan(
		&created.IssueID,
		&created.TeamID,
		&created.UserID,
		&created.UserName,
		&created.UserEmail,
		&created.IssueType,
		&created.Title,
		&created.Description,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return created, nil
}

func (r *IssueRepository) ListIssues(ctx context.Context) ([]models.Issue, error) {
	query := `
		SELECT issue_id, team_id, user_id, user_name, user_email,
		       issue_type, title, description, status, created_at, updated_at
		FROM issues
		ORDER BY created_at DESC
	`

	rows, err := r.databasePool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}
	defer rows.Close()

	results := make([]models.Issue, 0)
	for rows.Next() {
		var item models.Issue
		if err := rows.Scan(
			&item.IssueID,
			&item.TeamID,
			&item.UserID,
			&item.UserName,
			&item.UserEmail,
			&item.IssueType,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan issue: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate issues: %w", err)
	}

	return results, nil
}

func (r *IssueRepository) UpdateIssueStatus(ctx context.Context, issueID uuid.UUID, status string) error {
	query := `UPDATE issues SET status = $1, updated_at = now() WHERE issue_id = $2`
	result, err := r.databasePool.Exec(ctx, query, status, issueID)
	if err != nil {
		return fmt.Errorf("failed to update issue status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *IssueRepository) DeleteIssue(ctx context.Context, issueID uuid.UUID) error {
	query := `DELETE FROM issues WHERE issue_id = $1`
	result, err := r.databasePool.Exec(ctx, query, issueID)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
