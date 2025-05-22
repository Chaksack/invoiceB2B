package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"strconv"
	"time"
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

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (*dtos.LoginUserResponse, error)
	VerifyOTP(ctx context.Context, email, otp string) (*dtos.LoginUserResponse, error)
	RefreshToken(ctx context.Context, tokenStr string) (*dtos.RefreshTokenResponse, error)
	LogoutUser(ctx context.Context, tokenStr string) error
	Toggle2FA(ctx context.Context, userIDStr string, enable bool) error
	GetConfig() *config.Config
}

type authService struct {
	userRepo            repositories.UserRepository
	kycRepo             repositories.KYCRepository
	jwtService          JWTService
	emailService        EmailService
	otpService          OTPService
	notificationService NotificationService
	activityLogService  ActivityLogService
	cfg                 *config.Config
}

func NewAuthService(
	userRepo repositories.UserRepository,
	kycRepo repositories.KYCRepository,
	jwtService JWTService,
	emailService EmailService,
	otpService OTPService,
	notificationService NotificationService,
	activityLogService ActivityLogService,
	cfg *config.Config,
) AuthService { // Return the interface type
	return &authService{ // Return a pointer to the struct that implements the interface
		userRepo:            userRepo,
		kycRepo:             kycRepo,
		jwtService:          jwtService,
		emailService:        emailService,
		otpService:          otpService,
		notificationService: notificationService,
		activityLogService:  activityLogService,
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

	kycDetail := &models.KYCDetail{
		UserID: createdUser.ID,
		Status: models.KYCPending,
	}
	_, err = s.kycRepo.CreateOrUpdate(ctx, kycDetail)
	if err != nil {
		log.Printf("Error creating initial KYC record for user %d: %v", createdUser.ID, err)
	} else {
		createdUser.KYCID = &kycDetail.ID
		if _, err := s.userRepo.Update(ctx, createdUser); err != nil {
			log.Printf("Error updating user %d with KYCID: %v", createdUser.ID, err)
		}
	}

	go func() {
		subject := "Welcome to Invoice Financing App!"
		body := fmt.Sprintf("Hi %s,\n\nWelcome to our platform! Your account has been created successfully.\nPlease complete your KYC to start using our services.\n\nThanks,\nThe Team", createdUser.FirstName)
		err := s.emailService.SendEmail(createdUser.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send welcome email to %s: %v", createdUser.Email, err)
		}
	}()

	eventPayload := map[string]interface{}{"user_id": createdUser.ID, "email": createdUser.Email}
	if s.notificationService != nil {
		if err := s.notificationService.PublishUserRegisteredEvent(eventPayload); err != nil {
			log.Printf("Failed to publish user_registered event for %s: %v", createdUser.Email, err)
		}
	} else {
		log.Println("NotificationService is nil, skipping event publishing.")
	}

	_ = s.activityLogService.LogActivity(ctx, nil, &createdUser.ID, "USER_REGISTERED", fmt.Sprintf("User %s registered.", createdUser.Email), "")

	return createdUser, nil
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (*dtos.LoginUserResponse, error) {
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

	userResponse := dtos.UserResponse{ // Correctly map *models.User to dtos.UserResponse
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CompanyName:  user.CompanyName,
		IsActive:     user.IsActive,
		TwoFAEnabled: user.TwoFAEnabled,
	}

	if user.TwoFAEnabled {
		userIDStr := strconv.FormatUint(uint64(user.ID), 10)
		otp, err := s.otpService.GenerateAndStoreOTP(ctx, userIDStr)
		if err != nil {
			log.Printf("Failed to generate OTP for user %s: %v", user.Email, err)
			return nil, fmt.Errorf("failed to initiate 2FA: %w", err)
		}

		go func() {
			subject := "Your 2FA Login Code"
			body := fmt.Sprintf("Hi %s,\n\nYour One-Time Password for login is: %s\nIt will expire in %d minutes.\n\nThanks,\nThe Team", user.FirstName, otp, int(s.cfg.OTPExpirationMinutes.Minutes()))
			if emailErr := s.emailService.SendEmail(user.Email, subject, body); emailErr != nil {
				log.Printf("Failed to send 2FA OTP email to %s: %v", user.Email, emailErr)
			}
		}()

		return &dtos.LoginUserResponse{User: userResponse, TwoFARequired: true, Message: "OTP sent to your email for 2FA verification."}, nil
	}

	accessToken, accessExp, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	_ = s.activityLogService.LogActivity(ctx, nil, &user.ID, "USER_LOGIN_SUCCESS", fmt.Sprintf("User %s logged in successfully.", user.Email), "")

	return &dtos.LoginUserResponse{
		User:                 userResponse, // Use the mapped DTO
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		TwoFARequired:        false,
		AccessTokenExpiresAt: accessExp.Unix(),
		Message:              "Login successful.",
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, email, otp string) (*dtos.LoginUserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if !user.TwoFAEnabled {
		return nil, Err2FANotEnabled
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)
	valid, err := s.otpService.VerifyOTP(ctx, userIDStr, otp)
	if err != nil {
		log.Printf("Error verifying OTP for user %s: %v", user.Email, err)
		return nil, ErrOTPInvalidOrExpired
	}
	if !valid {
		_ = s.activityLogService.LogActivity(ctx, nil, &user.ID, "USER_LOGIN_2FA_FAILED", fmt.Sprintf("User %s failed 2FA OTP verification.", user.Email), "")
		return nil, ErrOTPInvalidOrExpired
	}

	accessToken, accessExp, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.otpService.DeleteOTP(ctx, userIDStr); err != nil {
		log.Printf("Warning: Failed to delete OTP for user %s after verification: %v", userIDStr, err)
	}

	userResponse := dtos.UserResponse{ // Correctly map *models.User to dtos.UserResponse
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CompanyName:  user.CompanyName,
		IsActive:     user.IsActive,
		TwoFAEnabled: user.TwoFAEnabled,
	}

	_ = s.activityLogService.LogActivity(ctx, nil, &user.ID, "USER_LOGIN_2FA_SUCCESS", fmt.Sprintf("User %s logged in successfully via 2FA.", user.Email), "")

	return &dtos.LoginUserResponse{
		User:                 userResponse, // Use the mapped DTO
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		TwoFARequired:        false,
		AccessTokenExpiresAt: accessExp.Unix(),
		Message:              "OTP verified successfully. Login complete.",
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, tokenStr string) (*dtos.RefreshTokenResponse, error) {
	claims, err := s.jwtService.ValidateToken(tokenStr, true)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	isBlacklisted, _ := s.otpService.IsTokenBlacklisted(ctx, tokenStr)
	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrRefreshTokenInvalid
	}
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}
	userID := uint(parsedUserID)

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

	newRefreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return &dtos.RefreshTokenResponse{
		AccessToken:          newAccessToken,
		RefreshToken:         newRefreshToken,
		AccessTokenExpiresAt: accessExp.Unix(),
		Message:              "Token refreshed successfully.",
	}, nil
}

func (s *authService) LogoutUser(ctx context.Context, tokenStr string) error {
	claims, err := s.jwtService.ValidateToken(tokenStr, false)
	if err != nil {
		log.Printf("Logout: Validating access token failed (possibly expired): %v", err)
	}

	var expiryDuration time.Duration
	if claims != nil {
		if expFloat, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(expFloat), 0)
			if expTime.After(time.Now()) {
				expiryDuration = time.Until(expTime)
			} else {
				expiryDuration = time.Minute
			}
		} else {
			expiryDuration = s.cfg.JWTAccessTokenExpirationMinutes
		}
	} else {
		expiryDuration = s.cfg.JWTAccessTokenExpirationMinutes
	}

	err = s.otpService.BlacklistToken(ctx, tokenStr, expiryDuration)
	if err != nil {
		log.Printf("Failed to blacklist access token on logout: %v", err)
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	if claims != nil {
		userIDStr, ok := claims["user_id"].(string)
		if ok {
			parsedUserID, pErr := strconv.ParseUint(userIDStr, 10, 64)
			if pErr == nil {
				uid := uint(parsedUserID)
				_ = s.activityLogService.LogActivity(ctx, nil, &uid, "USER_LOGOUT", fmt.Sprintf("User ID %d logged out.", uid), "")
			}
		}
	}
	return nil
}

func (s *authService) Toggle2FA(ctx context.Context, userIDStr string, enable bool) error {
	parsedUserID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}
	userID := uint(parsedUserID)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	user.TwoFAEnabled = enable
	if !enable {
		user.EmailOTP = nil
		user.EmailOTPExp = nil
	}

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		log.Printf("Failed to update 2FA status for user %d: %v", userID, err)
		return fmt.Errorf("could not update 2FA status: %w", err)
	}

	action := "USER_2FA_DISABLED"
	if enable {
		action = "USER_2FA_ENABLED"
	}
	_ = s.activityLogService.LogActivity(ctx, nil, &userID, action, fmt.Sprintf("User ID %d %s 2FA.", userID, action), "")

	return nil
}
