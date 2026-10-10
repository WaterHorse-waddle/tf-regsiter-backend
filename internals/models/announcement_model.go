package models

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	AnnouncementID uuid.UUID `json:"announcement_id" db:"announcement_id"`
	Title          string    `json:"title" db:"title"`
	Body           string    `json:"body" db:"body"`
	CreatedByEmail string    `json:"created_by_email,omitempty" db:"created_by_email"`
	CreatedByName  string    `json:"created_by_name,omitempty" db:"created_by_name"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAnnouncementRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}