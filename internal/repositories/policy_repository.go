package repositories

import (
	"context"
	"invoiceB2B/internal/models"
	"time"

	"gorm.io/gorm"
)

// PolicyRepository handles database operations for authorization policies
type PolicyRepository interface {
	// FindByID finds a policy by its ID
	FindByID(ctx context.Context, id uint) (*models.Policy, error)
	
	// FindByPolicyID finds a policy by its unique policy ID
	FindByPolicyID(ctx context.Context, policyID string) (*models.Policy, error)
	
	// FindAll returns all policies, optionally filtered by active status
	FindAll(ctx context.Context, activeOnly bool) ([]models.Policy, error)
	
	// Create creates a new policy
	Create(ctx context.Context, policy *models.Policy) (*models.Policy, error)
	
	// Update updates an existing policy
	Update(ctx context.Context, policy *models.Policy) (*models.Policy, error)
	
	// Delete deletes a policy
	Delete(ctx context.Context, id uint) error
	
	// LogAuthorizationDecision logs an authorization decision
	LogAuthorizationDecision(ctx context.Context, log *models.AuthorizationLog) error
	
	// GetAuthorizationLogs retrieves authorization logs with optional filtering
	GetAuthorizationLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.AuthorizationLog, error)
}

type policyRepository struct {
	db *gorm.DB
}

// NewPolicyRepository creates a new policy repository
func NewPolicyRepository(db *gorm.DB) PolicyRepository {
	return &policyRepository{db: db}
}

func (r *policyRepository) FindByID(ctx context.Context, id uint) (*models.Policy, error) {
	var policy models.Policy
	result := r.db.WithContext(ctx).First(&policy, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &policy, nil
}

func (r *policyRepository) FindByPolicyID(ctx context.Context, policyID string) (*models.Policy, error) {
	var policy models.Policy
	result := r.db.WithContext(ctx).Where("policy_id = ?", policyID).First(&policy)
	if result.Error != nil {
		return nil, result.Error
	}
	return &policy, nil
}

func (r *policyRepository) FindAll(ctx context.Context, activeOnly bool) ([]models.Policy, error) {
	var policies []models.Policy
	query := r.db.WithContext(ctx)
	
	if activeOnly {
		query = query.Where("active = ?", true)
	}
	
	// Order by priority (higher priority first)
	result := query.Order("priority DESC").Find(&policies)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return policies, nil
}

func (r *policyRepository) Create(ctx context.Context, policy *models.Policy) (*models.Policy, error) {
	result := r.db.WithContext(ctx).Create(policy)
	if result.Error != nil {
		return nil, result.Error
	}
	return policy, nil
}

func (r *policyRepository) Update(ctx context.Context, policy *models.Policy) (*models.Policy, error) {
	result := r.db.WithContext(ctx).Save(policy)
	if result.Error != nil {
		return nil, result.Error
	}
	return policy, nil
}

func (r *policyRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Policy{}, id)
	return result.Error
}

func (r *policyRepository) LogAuthorizationDecision(ctx context.Context, log *models.AuthorizationLog) error {
	// Ensure timestamp is set
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	
	result := r.db.WithContext(ctx).Create(log)
	return result.Error
}

func (r *policyRepository) GetAuthorizationLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.AuthorizationLog, error) {
	var logs []models.AuthorizationLog
	query := r.db.WithContext(ctx)
	
	// Apply filters
	for key, value := range filters {
		query = query.Where(key, value)
	}
	
	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	// Order by timestamp (newest first)
	result := query.Order("timestamp DESC").Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return logs, nil
}