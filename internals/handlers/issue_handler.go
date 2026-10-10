package handlers

import (
	"strings"

	"github.com/Anshul1310/tf-register/internals/models"
	"github.com/Anshul1310/tf-register/internals/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IssueHandler struct {
	issueService *service.IssueService
}

func NewIssueHandler(issueService *service.IssueService) *IssueHandler {
	return &IssueHandler{issueService: issueService}
}

// CreateIssue handles POST /issue — authenticated participant submits an issue.
func (h *IssueHandler) CreateIssue(c *fiber.Ctx) error {
	var req models.CreateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: "Invalid issue payload", Error: err.Error(),
		})
	}

	// Extract user info from middleware locals (set by AuthenticateUser middleware).
	// The middleware stores userID as uuid.UUID and userEmail as string.
	var userIDPtr *uuid.UUID
	if userIDVal := c.Locals("userID"); userIDVal != nil {
		if parsed, ok := userIDVal.(uuid.UUID); ok {
			userIDPtr = &parsed
		}
	}

	userEmail, _ := c.Locals("userEmail").(string)

	// teamID is optionally passed in the request body for context
	teamIDVal := strings.TrimSpace(req.TeamID)
	var teamIDPtr *string
	if teamIDVal != "" {
		teamIDPtr = &teamIDVal
	}

	issue, err := h.issueService.CreateIssue(
		c.Context(),
		req.IssueType,
		req.Title,
		req.Description,
		userIDPtr,
		req.UserName,
		userEmail,
		teamIDPtr,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: err.Error(), Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true, Message: "Issue submitted successfully", Data: issue,
	})
}

// ListIssues handles GET /admin/issues — admin-only list of all issues.
func (h *IssueHandler) ListIssues(c *fiber.Ctx) error {
	issues, err := h.issueService.ListIssues(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false, Message: "Failed to list issues", Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true, Message: "Issues retrieved successfully", Data: issues,
	})
}

// UpdateIssueStatus handles PATCH /admin/issues/:id/status — admin updates status.
func (h *IssueHandler) UpdateIssueStatus(c *fiber.Ctx) error {
	issueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: "Invalid issue ID",
		})
	}

	var req models.UpdateIssueStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: "Invalid status payload", Error: err.Error(),
		})
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: "Status is required",
		})
	}

	if err := h.issueService.UpdateIssueStatus(c.Context(), issueID, req.Status); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no rows") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false, Message: "Issue not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: err.Error(), Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true, Message: "Issue status updated",
	})
}

// DeleteIssue handles DELETE /admin/issues/:id — admin deletes an issue.
func (h *IssueHandler) DeleteIssue(c *fiber.Ctx) error {
	issueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false, Message: "Invalid issue ID",
		})
	}

	if err := h.issueService.DeleteIssue(c.Context(), issueID); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no rows") {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false, Message: "Issue not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false, Message: "Failed to delete issue", Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true, Message: "Issue deleted successfully",
	})
}
