package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"invoiceB2B/internal/models"
	"log"
)

// FinancialInstitutionRepository interface defines the operations for financial institution data management.
type FinancialInstitutionRepository interface {
	// Financial Institution methods
	CreateFinancialInstitution(ctx context.Context, fi *models.FinancialInstitution) error
	UpdateFinancialInstitution(ctx context.Context, fi *models.FinancialInstitution) error
	DeleteFinancialInstitution(ctx context.Context, id uint) error
	FindFinancialInstitutionByID(ctx context.Context, id uint) (*models.FinancialInstitution, error)
	FindFinancialInstitutionByCode(ctx context.Context, code string) (*models.FinancialInstitution, error)
	FindAllFinancialInstitutions(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.FinancialInstitution, int64, error)

	// Financial Institution Product methods
	CreateFinancialInstitutionProduct(ctx context.Context, product *models.FinancialInstitutionProduct) error
	UpdateFinancialInstitutionProduct(ctx context.Context, product *models.FinancialInstitutionProduct) error
	DeleteFinancialInstitutionProduct(ctx context.Context, id uint) error
	FindFinancialInstitutionProductByID(ctx context.Context, id uint) (*models.FinancialInstitutionProduct, error)
	FindFinancialInstitutionProductsByFinancialInstitutionID(ctx context.Context, fiID uint, page, pageSize int) ([]models.FinancialInstitutionProduct, int64, error)

	// Financial Institution Term methods
	CreateFinancialInstitutionTerm(ctx context.Context, term *models.FinancialInstitutionTerm) error
	UpdateFinancialInstitutionTerm(ctx context.Context, term *models.FinancialInstitutionTerm) error
	DeleteFinancialInstitutionTerm(ctx context.Context, id uint) error
	FindFinancialInstitutionTermByID(ctx context.Context, id uint) (*models.FinancialInstitutionTerm, error)
	FindFinancialInstitutionTermsByFinancialInstitutionID(ctx context.Context, fiID uint, page, pageSize int) ([]models.FinancialInstitutionTerm, int64, error)
	FindFinancialInstitutionTermsByProductID(ctx context.Context, productID uint, page, pageSize int) ([]models.FinancialInstitutionTerm, int64, error)
}

type financialInstitutionRepository struct {
	db *gorm.DB
}

// NewFinancialInstitutionRepository creates a new instance of FinancialInstitutionRepository.
func NewFinancialInstitutionRepository(db *gorm.DB) FinancialInstitutionRepository {
	return &financialInstitutionRepository{db: db}
}

// Financial Institution methods implementation

func (r *financialInstitutionRepository) CreateFinancialInstitution(ctx context.Context, fi *models.FinancialInstitution) error {
	if err := r.db.WithContext(ctx).Create(fi).Error; err != nil {
		log.Printf("Error creating financial institution in DB: %v", err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) UpdateFinancialInstitution(ctx context.Context, fi *models.FinancialInstitution) error {
	if err := r.db.WithContext(ctx).Save(fi).Error; err != nil {
		log.Printf("Error updating financial institution %d in DB: %v", fi.ID, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) DeleteFinancialInstitution(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.FinancialInstitution{}, id).Error; err != nil {
		log.Printf("Error deleting financial institution %d from DB: %v", id, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionByID(ctx context.Context, id uint) (*models.FinancialInstitution, error) {
	var fi models.FinancialInstitution
	if err := r.db.WithContext(ctx).First(&fi, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("financial institution not found")
		}
		log.Printf("Error finding financial institution by ID %d in DB: %v", id, err)
		return nil, err
	}
	return &fi, nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionByCode(ctx context.Context, code string) (*models.FinancialInstitution, error) {
	var fi models.FinancialInstitution
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&fi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("financial institution not found")
		}
		log.Printf("Error finding financial institution by code %s in DB: %v", code, err)
		return nil, err
	}
	return &fi, nil
}

func (r *financialInstitutionRepository) FindAllFinancialInstitutions(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.FinancialInstitution, int64, error) {
	var financialInstitutions []models.FinancialInstitution
	var total int64

	query := r.db.WithContext(ctx).Model(&models.FinancialInstitution{})

	// Apply filters
	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code, ok := filters["code"]; ok && code != "" {
		query = query.Where("code = ?", code)
	}
	if isActive, ok := filters["is_active"]; ok && isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting financial institutions with filters: %v", err)
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&financialInstitutions).Error; err != nil {
		log.Printf("Error fetching financial institutions with filters: %v", err)
		return nil, 0, err
	}

	return financialInstitutions, total, nil
}

// Financial Institution Product methods implementation

func (r *financialInstitutionRepository) CreateFinancialInstitutionProduct(ctx context.Context, product *models.FinancialInstitutionProduct) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		log.Printf("Error creating financial institution product in DB: %v", err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) UpdateFinancialInstitutionProduct(ctx context.Context, product *models.FinancialInstitutionProduct) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		log.Printf("Error updating financial institution product %d in DB: %v", product.ID, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) DeleteFinancialInstitutionProduct(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.FinancialInstitutionProduct{}, id).Error; err != nil {
		log.Printf("Error deleting financial institution product %d from DB: %v", id, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionProductByID(ctx context.Context, id uint) (*models.FinancialInstitutionProduct, error) {
	var product models.FinancialInstitutionProduct
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("financial institution product not found")
		}
		log.Printf("Error finding financial institution product by ID %d in DB: %v", id, err)
		return nil, err
	}
	return &product, nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionProductsByFinancialInstitutionID(ctx context.Context, fiID uint, page, pageSize int) ([]models.FinancialInstitutionProduct, int64, error) {
	var products []models.FinancialInstitutionProduct
	var total int64

	query := r.db.WithContext(ctx).Model(&models.FinancialInstitutionProduct{}).Where("financial_institution_id = ?", fiID)

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting financial institution products for financial institution %d: %v", fiID, err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		log.Printf("Error fetching financial institution products for financial institution %d: %v", fiID, err)
		return nil, 0, err
	}

	return products, total, nil
}

// Financial Institution Term methods implementation

func (r *financialInstitutionRepository) CreateFinancialInstitutionTerm(ctx context.Context, term *models.FinancialInstitutionTerm) error {
	if err := r.db.WithContext(ctx).Create(term).Error; err != nil {
		log.Printf("Error creating financial institution term in DB: %v", err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) UpdateFinancialInstitutionTerm(ctx context.Context, term *models.FinancialInstitutionTerm) error {
	if err := r.db.WithContext(ctx).Save(term).Error; err != nil {
		log.Printf("Error updating financial institution term %d in DB: %v", term.ID, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) DeleteFinancialInstitutionTerm(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.FinancialInstitutionTerm{}, id).Error; err != nil {
		log.Printf("Error deleting financial institution term %d from DB: %v", id, err)
		return err
	}
	return nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionTermByID(ctx context.Context, id uint) (*models.FinancialInstitutionTerm, error) {
	var term models.FinancialInstitutionTerm
	if err := r.db.WithContext(ctx).First(&term, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("financial institution term not found")
		}
		log.Printf("Error finding financial institution term by ID %d in DB: %v", id, err)
		return nil, err
	}
	return &term, nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionTermsByFinancialInstitutionID(ctx context.Context, fiID uint, page, pageSize int) ([]models.FinancialInstitutionTerm, int64, error) {
	var terms []models.FinancialInstitutionTerm
	var total int64

	query := r.db.WithContext(ctx).Model(&models.FinancialInstitutionTerm{}).Where("financial_institution_id = ?", fiID)

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting financial institution terms for financial institution %d: %v", fiID, err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&terms).Error; err != nil {
		log.Printf("Error fetching financial institution terms for financial institution %d: %v", fiID, err)
		return nil, 0, err
	}

	return terms, total, nil
}

func (r *financialInstitutionRepository) FindFinancialInstitutionTermsByProductID(ctx context.Context, productID uint, page, pageSize int) ([]models.FinancialInstitutionTerm, int64, error) {
	var terms []models.FinancialInstitutionTerm
	var total int64

	query := r.db.WithContext(ctx).Model(&models.FinancialInstitutionTerm{}).Where("financial_institution_product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting financial institution terms for product %d: %v", productID, err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&terms).Error; err != nil {
		log.Printf("Error fetching financial institution terms for product %d: %v", productID, err)
		return nil, 0, err
	}

	return terms, total, nil
}
