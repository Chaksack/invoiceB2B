package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailExists         = errors.New("user with this email already exists")
	ErrOTPInvalidOrExpired = errors.New("otp is invalid or has expired")
	Err2FANotEnabled       = errors.New("2fa is not enabled for this user")
	ErrAccountNotActive    = errors.New("user account is not active")
	ErrKYCNotApproved      = errors.New("user kyc not approved")
	ErrRefreshTokenInvalid = errors.New("refresh token is invalid or expired")
	ErrTokenBlacklisted    = errors.New("token has been blacklisted")
)

type LoginResult struct {
	User                 *models.User
	AccessToken          string
	RefreshToken         string
	TwoFARequired        bool
	AccessTokenExpiresAt int64
}

type RefreshTokenResult struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresAt int64
}

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (*LoginResult, error)
	VerifyOTP(ctx context.Context, email, otp string) (*LoginResult, error)
	RefreshToken(ctx context.Context, tokenStr string) (*RefreshTokenResult, error)
	LogoutUser(ctx context.Context, tokenStr string) error
	Toggle2FA(ctx context.Context, userID string, enable bool) error
	GetConfig() *config.Config // Helper to access config, e.g., for cookie settings in handler
}

type authService struct {
	userRepo            repositories.UserRepository
	kycRepo             repositories.KYCRepository
	jwtService          JWTService
	emailService        EmailService
	otpService          OTPService
	notificationService NotificationService // For RabbitMQ events
	cfg                 *config.Config
}

func NewAuthService(
	userRepo repositories.UserRepository,
	kycRepo repositories.KYCRepository,
	jwtService JWTService,
	emailService EmailService,
	otpService OTPService,
	notificationService NotificationService,
	cfg *config.Config,
) AuthService {
	return &authService{
		userRepo:            userRepo,
		kycRepo:             kycRepo,
		jwtService:          jwtService,
		emailService:        emailService,
		otpService:          otpService,
		notificationService: notificationService,
		cfg:                 cfg,
	}
}

func (s *authService) GetConfig() *config.Config {
	return s.cfg
}

func (s *authService) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("Error creating user in DB: %v", err)
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	// Create initial KYC record
	kycDetail := &models.KYCDetail{
		UserID: createdUser.ID,
		Status: models.KYCPending, // Default status
	}
	_, err = s.kycRepo.CreateOrUpdate(ctx, kycDetail)
	if err != nil {
		log.Printf("Error creating initial KYC record for user %s: %v", createdUser.ID, err)
		// Proceed with user creation, but log this issue. Might need a cleanup mechanism.
	} else {
		createdUser.KYCID = &kycDetail.ID
		if _, err := s.userRepo.Update(ctx, createdUser); err != nil {
			log.Printf("Error updating user %s with KYCID: %v", createdUser.ID, err)
		}
	}

	// Send welcome email (asynchronously via RabbitMQ ideally)
	go func() { // Fire and forget for now, proper worker setup needed for RabbitMQ
		subject := "Welcome to Invoice Financing App!"
		body := fmt.Sprintf("Hi %s,\n\nWelcome to our platform! Your account has been created successfully.\nPlease complete your KYC to start using our services.\n\nThanks,\nThe Team", createdUser.FirstName)
		err := s.emailService.SendEmail(createdUser.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send welcome email to %s: %v", createdUser.Email, err)
		}
	}()

	// Publish user_registered event (placeholder for RabbitMQ)
	eventPayload := map[string]interface{}{"user_id": createdUser.ID.String(), "email": createdUser.Email}
	if err := s.notificationService.PublishUserRegisteredEvent(eventPayload); err != nil {
		log.Printf("Failed to publish user_registered event for %s: %v", createdUser.Email, err)
	}

	return createdUser, nil
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if !models.CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountNotActive
	}

	// Check KYC status - for this app, KYC approval is needed to use features, not strictly for login.
	// We can add a check here if login itself should be blocked without KYC.
	// For now, let's assume login is allowed, but subsequent actions might be restricted.

	if user.TwoFAEnabled {
		otp, err := s.otpService.GenerateAndStoreOTP(ctx, user.ID.String())
		if err != nil {
			log.Printf("Failed to generate OTP for user %s: %v", user.Email, err)
			return nil, fmt.Errorf("failed to initiate 2FA: %w", err)
		}

		// Send OTP via email
		go func() { // Fire and forget
			subject := "Your 2FA Login Code"
			body := fmt.Sprintf("Hi %s,\n\nYour One-Time Password for login is: %s\nIt will expire in %d minutes.\n\nThanks,\nThe Team", user.FirstName, otp, int(s.cfg.OTPExpirationMinutes.Minutes()))
			if emailErr := s.emailService.SendEmail(user.Email, subject, body); emailErr != nil {
				log.Printf("Failed to send 2FA OTP email to %s: %v", user.Email, emailErr)
			}
		}()

		return &LoginResult{User: user, TwoFARequired: true}, nil
	}

	// No 2FA, proceed with token generation
	accessToken, accessExp, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResult{
		User:                 user,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		TwoFARequired:        false,
		AccessTokenExpiresAt: accessExp.Unix(),
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, email, otp string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if !user.TwoFAEnabled { // Should not happen if login flow is correct, but good check
		return nil, Err2FANotEnabled
	}

	valid, err := s.otpService.VerifyOTP(ctx, user.ID.String(), otp)
	if err != nil {
		log.Printf("Error verifying OTP for user %s: %v", user.Email, err)
		return nil, ErrOTPInvalidOrExpired
	}
	if !valid {
		return nil, ErrOTPInvalidOrExpired
	}

	// OTP is valid, generate tokens
	accessToken, accessExp, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear OTP after successful verification
	if err := s.otpService.DeleteOTP(ctx, user.ID.String()); err != nil {
		log.Printf("Warning: Failed to delete OTP for user %s after verification: %v", user.ID, err)
	}

	return &LoginResult{
		User:                 user,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		TwoFARequired:        false, // Verification successful
		AccessTokenExpiresAt: accessExp.Unix(),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, tokenStr string) (*RefreshTokenResult, error) {
	claims, err := s.jwtService.ValidateToken(tokenStr, true) // true for refresh token
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Check if token is blacklisted (for logout)
	isBlacklisted, _ := s.otpService.IsTokenBlacklisted(ctx, tokenStr)
	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrRefreshTokenInvalid
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrAccountNotActive
	}

	newAccessToken, accessExp, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Optionally, rotate refresh token
	newRefreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}
	// If rotating, you might want to blacklist the old refresh token
	// s.otpService.BlacklistToken(ctx, tokenStr, s.cfg.JWTRefreshTokenExpirationDays)

	return &RefreshTokenResult{
		AccessToken:          newAccessToken,
		RefreshToken:         newRefreshToken, // Send new refresh token
		AccessTokenExpiresAt: accessExp.Unix(),
	}, nil
}

func (s *authService) LogoutUser(ctx context.Context, tokenStr string) error {
	claims, err := s.jwtService.ValidateToken(tokenStr, false) // false for access token
	if err != nil {
		// If token is already expired, that's fine for logout.
		// If it's invalid for other reasons, log it but proceed.
		log.Printf("Logout: Validating access token failed (possibly expired): %v", err)
	}

	var expiryDuration time.Duration
	if claims != nil {
		if expFloat, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(expFloat), 0)
			if expTime.After(time.Now()) {
				expiryDuration = time.Until(expTime)
			} else {
				expiryDuration = time.Minute // Already expired, blacklist for a short while just in case
			}
		} else {
			expiryDuration = s.cfg.JWTAccessTokenExpirationMinutes // Default if 'exp' is not float64
		}
	} else {
		expiryDuration = s.cfg.JWTAccessTokenExpirationMinutes // Default if claims are nil
	}

	// Blacklist the access token until its original expiry
	err = s.otpService.BlacklistToken(ctx, tokenStr, expiryDuration)
	if err != nil {
		log.Printf("Failed to blacklist access token on logout: %v", err)
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	// Note: Refresh token invalidation is harder if not stored server-side.
	// Client should discard the refresh token.
	// If refresh tokens are one-time use or stored server-side, invalidate it here.
	return nil
}

func (s *authService) Toggle2FA(ctx context.Context, userIDStr string, enable bool) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	user.TwoFAEnabled = enable
	if !enable { // If disabling, clear any stored OTP info (though not strictly used for email OTP model)
		user.EmailOTP = nil
		user.EmailOTPExp = nil
	}

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		log.Printf("Failed to update 2FA status for user %s: %v", userID, err)
		return fmt.Errorf("could not update 2FA status: %w", err)
	}
	return nil
}
