package repositories

import (
	"context"
	"errors"
	"invoiceB2B/internal/models"
	"log"

	"gorm.io/gorm"
)

// UserRepository interface defines the operations for user data management.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByIDWithKYC(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	CountAll(ctx context.Context) (int64, error) // Added for analytics
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create persists a new user record to the database.
func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("Error creating user in DB: %v", err)
		return nil, err
	}
	return user, nil
}

// FindByEmail retrieves a user by their email address.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found, so service layer can decide error type
		}
		log.Printf("Error finding user by email %s in DB: %v", email, err)
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user by their ID.
func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") // Consistent error for not found
		}
		log.Printf("Error finding user by ID %d in DB: %v", id, err)
		return nil, err
	}
	return &user, nil
}

// FindByIDWithKYC retrieves a user by ID, preloading their KYC details.
func (r *userRepository) FindByIDWithKYC(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("KYCDetail").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") // Consistent error for not found
		}
		log.Printf("Error finding user by ID %d with KYC in DB: %v", id, err)
		return nil, err
	}
	return &user, nil
}

// Update saves changes to an existing user record.
func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Printf("Error updating user %d in DB: %v", user.ID, err)
		return nil, err
	}
	return user, nil
}

// FindAllWithPagination retrieves a paginated list of all users.
func (r *userRepository) FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting users: %v", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Preload("KYCDetail").Find(&users).Error; err != nil {
		log.Printf("Error fetching users with pagination: %v", err)
		return nil, 0, err
	}
	return users, total, nil
}

// CountAll returns the total number of users.
func (r *userRepository) CountAll(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("Error counting all users in DB: %v", err)
		return 0, err
	}
	return total, nil
}
