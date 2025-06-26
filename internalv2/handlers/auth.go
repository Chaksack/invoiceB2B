package handlers

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"invoiceB2B/internal@v2/config"
	"invoiceB2B/internal@v2/database"
	"invoiceB2B/internal@v2/middleware"
	"invoiceB2B/internal@v2/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db  *database.DB
	cfg *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Role == "" {
		req.Role = "business"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Check if user already exists
	var existingUser models.User
	err = h.db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID)
	if err == nil {
		return fiber.NewError(fiber.StatusConflict, "Email already registered")
	} else if err != sql.ErrNoRows {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Insert new user
	var user models.User
	isApproved := req.Role == "admin"
	err = h.db.QueryRow(
		"INSERT INTO users (email, password_hash, role, is_approved) VALUES ($1, $2, $3, $4) RETURNING id, email, role, is_approved, created_at",
		req.Email, hashedPassword, req.Role, isApproved,
	).Scan(&user.ID, &user.Email, &user.Role, &user.IsApproved, &user.CreatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully. Awaiting admin approval.",
		"data": models.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Role:       user.Role,
			IsApproved: user.IsApproved,
			CreatedAt:  user.CreatedAt,
		},
		"timestamp": fiber.Now(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get user from database
	var user models.User
	err := h.db.QueryRow(
		"SELECT id, email, password_hash, role, is_approved, created_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsApproved, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if business user is approved
	if !user.IsApproved && user.Role == "business" {
		return fiber.NewError(fiber.StatusForbidden, "Account awaiting admin approval")
	}

	// Generate JWT token
	claims := middleware.Claims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.cfg.JWTAccessTokenExpirationMinutes)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": models.LoginResponse{
			Token: tokenString,
			User: models.UserResponse{
				ID:         user.ID,
				Email:      user.Email,
				Role:       user.Role,
				IsApproved: user.IsApproved,
				CreatedAt:  user.CreatedAt,
			},
		},
		"timestamp": fiber.Now(),
	})
}

// Profile handles getting user profile
func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)

	// Get user from database
	var userData models.User
	err := h.db.QueryRow(
		"SELECT id, email, role, is_approved, created_at FROM users WHERE id = $1",
		user.ID,
	).Scan(&userData.ID, &userData.Email, &userData.Role, &userData.IsApproved, &userData.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": models.UserResponse{
			ID:         userData.ID,
			Email:      userData.Email,
			Role:       userData.Role,
			IsApproved: userData.IsApproved,
			CreatedAt:  userData.CreatedAt,
		},
		"timestamp": fiber.Now(),
	})
} 