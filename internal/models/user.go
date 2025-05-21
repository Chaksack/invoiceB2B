package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"` // Default handled by BeforeCreate
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	FirstName    string    `gorm:"type:varchar(50);not null"`
	LastName     string    `gorm:"type:varchar(50);not null"`
	CompanyName  string    `gorm:"type:varchar(100);not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`

	KYCID     *uuid.UUID `gorm:"type:uuid;null"`
	KYCDetail *KYCDetail `gorm:"foreignKey:KYCID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsActive     bool       `gorm:"default:true"`
	TwoFASecret  *string    `gorm:"type:varchar(255);null"` // For TOTP apps, not used for email OTP
	TwoFAEnabled bool       `gorm:"default:false"`          // General flag for 2FA
	EmailOTP     *string    `gorm:"type:varchar(10);null"`  // Store OTP for email 2FA
	EmailOTPExp  *time.Time `gorm:"null"`                   // Expiry for email OTP

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Invoices []Invoice `gorm:"foreignKey:UserID"`
}

type KYCStatus string

const (
	KYCPending          KYCStatus = "pending"
	KYCApproved         KYCStatus = "approved"
	KYCRejected         KYCStatus = "rejected"
	KYCResubmitRequired KYCStatus = "resubmit_required"
)

type KYCDetail struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID          uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Status          KYCStatus `gorm:"type:varchar(20);default:'pending';not null"`
	SubmittedAt     *time.Time
	ReviewedByID    *uuid.UUID `gorm:"type:uuid;null"`
	ReviewedBy      *Staff     `gorm:"foreignKey:ReviewedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReviewedAt      *time.Time
	RejectionReason *string   `gorm:"type:text;null"`
	DocumentsInfo   string    `gorm:"type:jsonb;null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	FirstName    string    `gorm:"type:varchar(50);not null"`
	LastName     string    `gorm:"type:varchar(50);not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:varchar(50);not null"`
	IsActive     bool      `gorm:"default:true"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.PasswordHash != "" {
		hashedPassword, err := HashPassword(u.PasswordHash)
		if err != nil {
			return err
		}
		u.PasswordHash = hashedPassword
	}
	return
}

func (s *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.PasswordHash != "" {
		hashedPassword, err := HashPassword(s.PasswordHash)
		if err != nil {
			return err
		}
		s.PasswordHash = hashedPassword
	}
	return
}

func (kyc *KYCDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if kyc.ID == uuid.Nil {
		kyc.ID = uuid.New()
	}
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
