package services

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dtos.UpdateUserProfileRequest) (*models.User, error)
	SubmitOrUpdateKYC(ctx context.Context, userID uuid.UUID, req dtos.SubmitKYCRequest) (*models.KYCDetail, error)
	GetKYCStatus(ctx context.Context, userID uuid.UUID) (*models.KYCDetail, error)
}

type userService struct {
	userRepo repositories.UserRepository
	kycRepo  repositories.KYCRepository
}

func NewUserService(userRepo repositories.UserRepository, kycRepo repositories.KYCRepository) UserService {
	return &userService{
		userRepo: userRepo,
		kycRepo:  kycRepo,
	}
}

func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByIDWithKYC(ctx, userID) // Assuming FindByIDWithKYC preloads KYCDetail
	if err != nil {
		log.Printf("Error fetching profile for user %s: %v", userID, err)
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dtos.UpdateUserProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.CompanyName != "" {
		user.CompanyName = req.CompanyName
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		log.Printf("Error updating profile for user %s: %v", userID, err)
		return nil, fmt.Errorf("could not update user profile: %w", err)
	}

	// Re-fetch with KYC to ensure the response DTO is complete
	fullUser, err := s.userRepo.FindByIDWithKYC(ctx, userID)
	if err != nil {
		log.Printf("Error re-fetching profile with KYC for user %s after update: %v", userID, err)
		return updatedUser, nil // Return partially updated user if re-fetch fails
	}
	return fullUser, nil
}

func (s *userService) SubmitOrUpdateKYC(ctx context.Context, userID uuid.UUID, req dtos.SubmitKYCRequest) (*models.KYCDetail, error) {
	_, err := s.userRepo.FindByID(ctx, userID) // Ensure user exists
	if err != nil {
		return nil, ErrUserNotFound
	}

	kycDetail, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Changed from string comparison
		log.Printf("Error finding KYC for user %s: %v", userID, err)
		return nil, fmt.Errorf("could not retrieve existing KYC: %w", err)
	}

	if kycDetail == nil { // No existing KYC, create new
		kycDetail = &models.KYCDetail{
			UserID: userID,
		}
	}

	// Update KYC details
	// In a real app, map fields from req to kycDetail carefully
	// For example:
	// kycDetail.BusinessRegistrationNumber = req.BusinessRegistrationNumber (if this field exists on model)
	kycDetail.DocumentsInfo = req.DocumentsInfo // Assuming DocumentsInfo is a JSON string of metadata
	kycDetail.Status = models.KYCPending        // Set status to pending on new submission/update
	now := time.Now()
	kycDetail.SubmittedAt = &now
	kycDetail.ReviewedAt = nil // Clear previous review if any
	kycDetail.ReviewedByID = nil
	kycDetail.RejectionReason = nil

	updatedKYC, err := s.kycRepo.CreateOrUpdate(ctx, kycDetail)
	if err != nil {
		log.Printf("Error creating/updating KYC for user %s: %v", userID, err)
		return nil, fmt.Errorf("could not save KYC information: %w", err)
	}
	return updatedKYC, nil
}

func (s *userService) GetKYCStatus(ctx context.Context, userID uuid.UUID) (*models.KYCDetail, error) {
	kycDetail, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Changed from string comparison
			return nil, errors.New("KYC record not found for this user") // Specific error for not found
		}
		log.Printf("Error retrieving KYC status for user %s: %v", userID, err)
		return nil, fmt.Errorf("could not retrieve KYC status: %w", err)
	}
	return kycDetail, nil
}
