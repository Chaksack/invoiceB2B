package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID uint, req dtos.UpdateUserProfileRequest) (*models.User, error)
	SubmitOrUpdateKYC(ctx context.Context, userID uint, req dtos.SubmitKYCRequest) (*models.KYCDetail, error)
	GetKYCStatus(ctx context.Context, userID uint) (*models.KYCDetail, error)
}

type userService struct {
	userRepo           repositories.UserRepository
	kycRepo            repositories.KYCRepository
	activityLogService ActivityLogService
}

func NewUserService(userRepo repositories.UserRepository, kycRepo repositories.KYCRepository, activityLogService ActivityLogService) UserService {
	return &userService{
		userRepo:           userRepo,
		kycRepo:            kycRepo,
		activityLogService: activityLogService, // Added
	}
}

func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByIDWithKYC(ctx, userID)
	if err != nil {
		log.Printf("Error fetching profile for user %d: %v", userID, err)
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID uint, req dtos.UpdateUserProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// Store old values for logging if needed
	// oldFirstName := user.FirstName
	// ...

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
		log.Printf("Error updating profile for user %d: %v", userID, err)
		return nil, fmt.Errorf("could not update user profile: %w", err)
	}

	_ = s.activityLogService.LogActivity(ctx, nil, &userID, "USER_PROFILE_UPDATED", fmt.Sprintf("User ID %d updated their profile.", userID), "")

	fullUser, err := s.userRepo.FindByIDWithKYC(ctx, userID)
	if err != nil {
		log.Printf("Error re-fetching profile with KYC for user %d after update: %v", userID, err)
		return updatedUser, nil
	}
	return fullUser, nil
}

func (s *userService) SubmitOrUpdateKYC(ctx context.Context, userID uint, req dtos.SubmitKYCRequest) (*models.KYCDetail, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("SubmitOrUpdateKYC: User not found for ID %d: %v", userID, err)
		return nil, ErrUserNotFound
	}

	kycDetail, err := s.kycRepo.FindByUserID(ctx, userID)
	isNewSubmission := false
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isNewSubmission = true
			kycDetail = &models.KYCDetail{
				UserID: userID,
			}
		} else {
			log.Printf("SubmitOrUpdateKYC: Error finding KYC for user %d: %v", userID, err)
			return nil, fmt.Errorf("could not retrieve existing KYC details: %w", err)
		}
	}

	kycDetail.DocumentsInfo = req.DocumentsInfo
	kycDetail.Status = models.KYCPending
	now := time.Now()
	kycDetail.SubmittedAt = &now
	kycDetail.ReviewedAt = nil
	kycDetail.ReviewedByID = nil
	kycDetail.RejectionReason = nil

	updatedKYC, err := s.kycRepo.CreateOrUpdate(ctx, kycDetail)
	if err != nil {
		log.Printf("SubmitOrUpdateKYC: Error creating/updating KYC for user %d: %v", userID, err)
		return nil, fmt.Errorf("could not save KYC information: %w", err)
	}

	if user.KYCID == nil || (updatedKYC != nil && *user.KYCID != updatedKYC.ID) {
		user.KYCID = &updatedKYC.ID
		if _, updateErr := s.userRepo.Update(ctx, user); updateErr != nil {
			log.Printf("SubmitOrUpdateKYC: CRITICAL - Failed to update User.KYCID for UserID %d with KYCDetail.ID %d: %v", userID, updatedKYC.ID, updateErr)
		} else {
			log.Printf("SubmitOrUpdateKYC: Successfully updated User.KYCID for UserID %d to %d", userID, updatedKYC.ID)
		}
	}

	// 6. Log the activity.
	action := "USER_KYC_UPDATED"
	if isNewSubmission {
		action = "USER_KYC_SUBMITTED"
	}
	if s.activityLogService != nil {
		_ = s.activityLogService.LogActivity(ctx, nil, &userID, action,
			map[string]interface{}{"kyc_id": updatedKYC.ID, "status": updatedKYC.Status},
			"")
	}

	return updatedKYC, nil
}

func (s *userService) GetKYCStatus(ctx context.Context, userID uint) (*models.KYCDetail, error) {
	kycDetail, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("KYC record not found for this user")
		}
		log.Printf("Error retrieving KYC status for user %d: %v", userID, err)
		return nil, fmt.Errorf("could not retrieve KYC status: %w", err)
	}
	return kycDetail, nil
}
