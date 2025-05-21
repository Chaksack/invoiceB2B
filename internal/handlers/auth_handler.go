package handlers

import (
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	authService services.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService services.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dtos.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	user := &models.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		CompanyName:  req.CompanyName,
		PasswordHash: req.Password, // Hashed by BeforeCreate hook
	}

	createdUser, err := h.authService.RegisterUser(c.Context(), user)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to register user", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dtos.RegisterUserResponse{
		User: dtos.UserResponse{
			ID:           createdUser.ID,
			Email:        createdUser.Email,
			FirstName:    createdUser.FirstName,
			LastName:     createdUser.LastName,
			CompanyName:  createdUser.CompanyName,
			IsActive:     createdUser.IsActive,
			TwoFAEnabled: createdUser.TwoFAEnabled,
		},
		Message: "User registered successfully. A KYC record has been initiated.",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dtos.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	result, err := h.authService.LoginUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Login failed", err) // More specific errors handled in service
	}

	if result.TwoFARequired {
		return c.Status(fiber.StatusOK).JSON(dtos.LoginUserResponse{
			Message:       "OTP sent to your email for 2FA verification.",
			TwoFARequired: true,
			User: dtos.UserResponse{ // Send some user info even if 2FA is pending
				ID:    result.User.ID,
				Email: result.User.Email,
			},
		})
	}

	// Set refresh token in HTTPOnly cookie
	// c.Cookie(&fiber.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    result.RefreshToken,
	// 	Expires:  time.Now().Add(h.authService.GetConfig().JWTRefreshTokenExpirationDays),
	// 	HTTPOnly: true,
	// 	Secure:   h.authService.GetConfig().AppEnv == "production", // Secure in production
	// 	SameSite: "Lax",
	// })

	return c.Status(fiber.StatusOK).JSON(dtos.LoginUserResponse{
		User: dtos.UserResponse{
			ID:           result.User.ID,
			Email:        result.User.Email,
			FirstName:    result.User.FirstName,
			LastName:     result.User.LastName,
			CompanyName:  result.User.CompanyName,
			IsActive:     result.User.IsActive,
			TwoFAEnabled: result.User.TwoFAEnabled,
		},
		AccessToken:          result.AccessToken,
		RefreshToken:         result.RefreshToken, // Client should store this securely (e.g., localStorage)
		Message:              "Login successful.",
		TwoFARequired:        false,
		AccessTokenExpiresAt: result.AccessTokenExpiresAt,
	})
}

func (h *AuthHandler) Verify2FA(c *fiber.Ctx) error {
	var req dtos.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	result, err := h.authService.VerifyOTP(c.Context(), req.Email, req.OTP)
	if err != nil {
		return utils.HandleError(c, fiber.StatusUnauthorized, "OTP verification failed", err)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.VerifyOTPResponse{
		User: dtos.UserResponse{
			ID:           result.User.ID,
			Email:        result.User.Email,
			FirstName:    result.User.FirstName,
			LastName:     result.User.LastName,
			CompanyName:  result.User.CompanyName,
			IsActive:     result.User.IsActive,
			TwoFAEnabled: result.User.TwoFAEnabled,
		},
		AccessToken:          result.AccessToken,
		RefreshToken:         result.RefreshToken,
		Message:              "OTP verified successfully. Login complete.",
		AccessTokenExpiresAt: result.AccessTokenExpiresAt,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dtos.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	newTokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Failed to refresh token", err)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		AccessToken:          newTokens.AccessToken,
		RefreshToken:         newTokens.RefreshToken,
		AccessTokenExpiresAt: newTokens.AccessTokenExpiresAt,
		Message:              "Token refreshed successfully.",
	})
}

func (h *AuthHandler) Enable2FA(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req dtos.Enable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	// No specific validation for a single boolean here, but could be added.

	err := h.authService.Toggle2FA(c.Context(), userID, req.Enable)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update 2FA status", err)
	}

	message := "2FA disabled successfully."
	if req.Enable {
		message = "2FA enabled successfully. Future logins will require OTP."
	}

	return c.Status(fiber.StatusOK).JSON(dtos.Enable2FAResponse{
		Message:      message,
		TwoFAEnabled: req.Enable,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		// This case should ideally not happen if middleware is correctly applied
		// but as a fallback, we can try to get token from header if needed.
		// For now, assume middleware populates it.
		log.Println("Logout: No user claims found in context. This might indicate an issue with middleware setup or an unauthenticated request reaching a protected route.")
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated or token missing in context", nil)
	}

	tokenStr := claims.Raw // Get the raw token string

	err := h.authService.LogoutUser(c.Context(), tokenStr)
	if err != nil {
		// Log the error but still return success to client as logout is best-effort
		log.Printf("Error during token invalidation on logout: %v", err)
	}

	// Optionally clear refresh token cookie if it was set
	// c.Cookie(&fiber.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    "",
	// 	Expires:  time.Now().Add(-time.Hour), // Expire immediately
	// 	HTTPOnly: true,
	// 	Secure:   h.authService.GetConfig().AppEnv == "production",
	// 	SameSite: "Lax",
	// })

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
