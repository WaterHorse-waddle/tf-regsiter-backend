package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anshul1310/tf-register/internals/models"
	"github.com/Anshul1310/tf-register/internals/repository"
	"github.com/google/uuid"
)

var validIssueTypes = map[string]bool{
	"Network Connectivity and Access":       true,
	"Presentation Submission and Upload":    true,
	"Lounging":                              true,
	"Attendance and Academic Records":       true,
	"ID Card Distribution Coordination":     true,
	"Other":                                 true,
}

type IssueService struct {
	issueRepo *repository.IssueRepository
}

func NewIssueService(issueRepo *repository.IssueRepository) *IssueService {
	return &IssueService{issueRepo: issueRepo}
}

func (s *IssueService) CreateIssue(
	ctx context.Context,
	issueType, title, description string,
	userID *uuid.UUID,
	userName, userEmail string,
	teamID *string,
) (*models.Issue, error) {
	issueType = strings.TrimSpace(issueType)
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if issueType == "" {
		return nil, errors.New("issue type is required")
	}
	// Also allow "Other: <custom text>" submitted by the frontend
	isValidType := validIssueTypes[issueType] ||
		(strings.HasPrefix(issueType, "Other: ") && len(strings.TrimPrefix(issueType, "Other: ")) > 0)
	if !isValidType {
		return nil, fmt.Errorf("invalid issue type: %s", issueType)
	}
	if title == "" {
		return nil, errors.New("issue title is required")
	}
	if description == "" {
		return nil, errors.New("issue description is required")
	}

	issue := &models.Issue{
		IssueID:     uuid.New(),
		IssueType:   issueType,
		Title:       title,
		Description: description,
		Status:      "open",
		TeamID:      teamID,
		UserID:      userID,
	}

	if userName != "" {
		issue.UserName = &userName
	}
	if userEmail != "" {
		issue.UserEmail = &userEmail
	}

	return s.issueRepo.CreateIssue(ctx, issue)
}

func (s *IssueService) ListIssues(ctx context.Context) ([]models.Issue, error) {
	return s.issueRepo.ListIssues(ctx)
}

func (s *IssueService) UpdateIssueStatus(ctx context.Context, issueID uuid.UUID, status string) error {
	status = strings.TrimSpace(strings.ToLower(status))
	validStatuses := map[string]bool{
		"open":        true,
		"in_progress": true,
		"resolved":    true,
		"closed":      true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.issueRepo.UpdateIssueStatus(ctx, issueID, status)
}

func (s *IssueService) DeleteIssue(ctx context.Context, issueID uuid.UUID) error {
	if issueID == uuid.Nil {
		return errors.New("issue id is required")
	}
	if err := s.issueRepo.DeleteIssue(ctx, issueID); err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}
