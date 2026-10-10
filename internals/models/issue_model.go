package models

import (
	"time"

	"github.com/google/uuid"
)

type Issue struct {
	IssueID     uuid.UUID `json:"issue_id"              db:"issue_id"`
	TeamID      *string   `json:"team_id,omitempty"     db:"team_id"`
	UserID      *uuid.UUID `json:"user_id,omitempty"    db:"user_id"`
	UserName    *string   `json:"user_name,omitempty"   db:"user_name"`
	UserEmail   *string   `json:"user_email,omitempty"  db:"user_email"`
	IssueType   string    `json:"issue_type"            db:"issue_type"`
	Title       string    `json:"title"                 db:"title"`
	Description string    `json:"description"           db:"description"`
	Status      string    `json:"status"                db:"status"`
	CreatedAt   time.Time `json:"created_at"            db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"            db:"updated_at"`
}

type CreateIssueRequest struct {
	IssueType   string `json:"issue_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserName    string `json:"user_name,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
}

type UpdateIssueStatusRequest struct {
	Status string `json:"status"` // "open", "in_progress", "resolved", "closed"
}
