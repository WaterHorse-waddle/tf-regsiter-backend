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

type AnnouncementService struct {
	announcementRepo *repository.AnnouncementRepository
}

func NewAnnouncementService(announcementRepo *repository.AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepo: announcementRepo}
}

func (s *AnnouncementService) ListAnnouncements(ctx context.Context) ([]models.Announcement, error) {
	return s.announcementRepo.ListAnnouncements(ctx, true)
}

func (s *AnnouncementService) ListAllAnnouncements(ctx context.Context) ([]models.Announcement, error) {
	return s.announcementRepo.ListAnnouncements(ctx, false)
}

func (s *AnnouncementService) CreateAnnouncement(ctx context.Context, title, body, createdByEmail, createdByName string) (*models.Announcement, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if title == "" {
		return nil, errors.New("announcement title is required")
	}
	if body == "" {
		return nil, errors.New("announcement body is required")
	}

	return s.announcementRepo.CreateAnnouncement(ctx, &models.Announcement{
		AnnouncementID: uuid.New(),
		Title:          title,
		Body:           body,
		CreatedByEmail: createdByEmail,
		CreatedByName:  createdByName,
		IsActive:       true,
	})
}

func (s *AnnouncementService) DeleteAnnouncement(ctx context.Context, announcementID uuid.UUID) error {
	if announcementID == uuid.Nil {
		return errors.New("announcement id is required")
	}
	if err := s.announcementRepo.DeleteAnnouncement(ctx, announcementID); err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}
	return nil
}