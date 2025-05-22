package handlers

import (
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService services.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)

	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format in token", err)
	}
	userID := uint(parsedUserID)

	user, err := h.userService.GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "User profile not found", err)
	}

	kycStatus := "Not Submitted"
	if user.KYCDetail != nil {
		kycStatus = string(user.KYCDetail.Status)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.UserProfileResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CompanyName:  user.CompanyName,
		IsActive:     user.IsActive,
		TwoFAEnabled: user.TwoFAEnabled,
		KYCStatus:    kycStatus,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
}

func (h *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format in token", err)
	}
	userID := uint(parsedUserID)

	var req dtos.UpdateUserProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	updatedUser, err := h.userService.UpdateUserProfile(c.Context(), userID, req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update user profile", err)
	}

	kycStatus := "Not Submitted"
	if updatedUser.KYCDetail != nil {
		kycStatus = string(updatedUser.KYCDetail.Status)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.UserProfileResponse{
		ID:           updatedUser.ID,
		Email:        updatedUser.Email,
		FirstName:    updatedUser.FirstName,
		LastName:     updatedUser.LastName,
		CompanyName:  updatedUser.CompanyName,
		IsActive:     updatedUser.IsActive,
		TwoFAEnabled: updatedUser.TwoFAEnabled,
		KYCStatus:    kycStatus,
		CreatedAt:    updatedUser.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
	})
}

func (h *UserHandler) SubmitKYC(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format in token", err)
	}
	userID := uint(parsedUserID)

	var req dtos.SubmitKYCRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	kycDetail, err := h.userService.SubmitOrUpdateKYC(c.Context(), userID, req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to submit KYC information", err)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.KYCStatusResponse{
		UserID:          kycDetail.UserID,
		Status:          string(kycDetail.Status),
		SubmittedAt:     kycDetail.SubmittedAt,
		RejectionReason: kycDetail.RejectionReason,
		Message:         "KYC information submitted successfully. It is pending review.",
	})
}

func (h *UserHandler) GetKYCStatus(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format in token", err)
	}
	userID := uint(parsedUserID)

	kycDetail, err := h.userService.GetKYCStatus(c.Context(), userID)
	if err != nil {
		if err.Error() == "KYC record not found for this user" {
			return c.Status(fiber.StatusOK).JSON(dtos.KYCStatusResponse{
				UserID:  userID,
				Status:  "Not Submitted",
				Message: "KYC information has not been submitted yet.",
			})
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve KYC status", err)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.KYCStatusResponse{
		UserID:          kycDetail.UserID,
		Status:          string(kycDetail.Status),
		SubmittedAt:     kycDetail.SubmittedAt,
		ReviewedAt:      kycDetail.ReviewedAt,
		RejectionReason: kycDetail.RejectionReason,
		Message:         "KYC status retrieved successfully.",
	})
}
