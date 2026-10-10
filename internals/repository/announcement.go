package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anshul1310/tf-register/internals/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnnouncementRepository struct {
	databasePool *pgxpool.Pool
}

func NewAnnouncementRepository(databasePool *pgxpool.Pool) *AnnouncementRepository {
	return &AnnouncementRepository{databasePool: databasePool}
}

func (r *AnnouncementRepository) CreateAnnouncement(ctx context.Context, announcement *models.Announcement) (*models.Announcement, error) {
	if announcement == nil {
		return nil, errors.New("announcement cannot be nil")
	}
	if announcement.AnnouncementID == uuid.Nil {
		announcement.AnnouncementID = uuid.New()
	}

	query := `
		INSERT INTO announcements (
			announcement_id, title, body, created_by_email, created_by_name, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, COALESCE($6, true), now(), now())
		RETURNING announcement_id, title, body, created_by_email, created_by_name, is_active, created_at, updated_at;
	`

	created := &models.Announcement{}
	err := r.databasePool.QueryRow(
		ctx,
		query,
		announcement.AnnouncementID,
		strings.TrimSpace(announcement.Title),
		strings.TrimSpace(announcement.Body),
		strings.TrimSpace(announcement.CreatedByEmail),
		strings.TrimSpace(announcement.CreatedByName),
		announcement.IsActive,
	).Scan(
		&created.AnnouncementID,
		&created.Title,
		&created.Body,
		&created.CreatedByEmail,
		&created.CreatedByName,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	return created, nil
}

func (r *AnnouncementRepository) ListAnnouncements(ctx context.Context, activeOnly bool) ([]models.Announcement, error) {
	query := `
		SELECT announcement_id, title, body, created_by_email, created_by_name, is_active, created_at, updated_at
		FROM announcements
	`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.databasePool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list announcements: %w", err)
	}
	defer rows.Close()

	results := make([]models.Announcement, 0)
	for rows.Next() {
		var item models.Announcement
		if err := rows.Scan(
			&item.AnnouncementID,
			&item.Title,
			&item.Body,
			&item.CreatedByEmail,
			&item.CreatedByName,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan announcement: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate announcements: %w", err)
	}

	return results, nil
}

func (r *AnnouncementRepository) DeleteAnnouncement(ctx context.Context, announcementID uuid.UUID) error {
	query := `DELETE FROM announcements WHERE announcement_id = $1`
	result, err := r.databasePool.Exec(ctx, query, announcementID)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}