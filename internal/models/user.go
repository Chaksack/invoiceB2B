package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	FirstName    string `gorm:"type:varchar(50);not null"`
	LastName     string `gorm:"type:varchar(50);not null"`
	CompanyName  string `gorm:"type:varchar(100);not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`

	KYCID     *uint      `gorm:"null"`
	KYCDetail *KYCDetail `gorm:"foreignKey:KYCID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsActive     bool       `gorm:"default:true"`
	TwoFASecret  *string    `gorm:"type:varchar(255);null"`
	TwoFAEnabled bool       `gorm:"default:false"`
	EmailOTP     *string    `gorm:"type:varchar(10);null"`
	EmailOTPExp  *time.Time `gorm:"null"`

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
	gorm.Model
	UserID          uint      `gorm:"uniqueIndex;not null"`
	Status          KYCStatus `gorm:"type:varchar(20);default:'pending';not null"`
	SubmittedAt     *time.Time
	ReviewedByID    *uint  `gorm:"null"`
	ReviewedBy      *Staff `gorm:"foreignKey:ReviewedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReviewedAt      *time.Time
	RejectionReason *string `gorm:"type:text;null"`
	DocumentsInfo   string  `gorm:"type:jsonb;null"`
}

type Staff struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	FirstName    string `gorm:"type:varchar(50);not null"`
	LastName     string `gorm:"type:varchar(50);not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(50);not null"`
	IsActive     bool   `gorm:"default:true"`
	LastLoginAt  *time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
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
	if s.PasswordHash != "" {
		hashedPassword, err := HashPassword(s.PasswordHash)
		if err != nil {
			return err
		}
		s.PasswordHash = hashedPassword
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
