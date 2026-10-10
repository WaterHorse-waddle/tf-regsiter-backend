package handlers

import (
	"strings"

	"github.com/Anshul1310/tf-register/internals/models"
	"github.com/Anshul1310/tf-register/internals/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AnnouncementHandler struct {
	announcementService *service.AnnouncementService
}

func NewAnnouncementHandler(announcementService *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{announcementService: announcementService}
}

func (h *AnnouncementHandler) ListPublicAnnouncements(c *fiber.Ctx) error {
	announcements, err := h.announcementService.ListAnnouncements(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "Failed to load announcements", Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "Announcements retrieved successfully", Data: announcements})
}

func (h *AnnouncementHandler) ListAdminAnnouncements(c *fiber.Ctx) error {
	announcements, err := h.announcementService.ListAllAnnouncements(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "Failed to load announcements", Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "Announcements retrieved successfully", Data: announcements})
}

func (h *AnnouncementHandler) CreateAnnouncement(c *fiber.Ctx) error {
	var req models.CreateAnnouncementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "Invalid announcement payload", Error: err.Error()})
	}

	createdByEmail, _ := c.Locals("adminEmail").(string)
	createdByName, _ := c.Locals("adminName").(string)

	announcement, err := h.announcementService.CreateAnnouncement(c.Context(), req.Title, req.Body, createdByEmail, createdByName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: err.Error(), Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{Success: true, Message: "Announcement published successfully", Data: announcement})
}

func (h *AnnouncementHandler) DeleteAnnouncement(c *fiber.Ctx) error {
	announcementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "Invalid announcement ID"})
	}

	if err := h.announcementService.DeleteAnnouncement(c.Context(), announcementID); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no rows") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "Announcement not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "Failed to delete announcement", Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "Announcement deleted successfully"})
}