package models

import (
	"gorm.io/gorm"
	"time" // Keep for explicit fields if needed, but gorm.Model handles CreatedAt, UpdatedAt
)

type ActivityLog struct {
	// gorm.Model // Original line
	ID        uint           `gorm:"primarykey"` // Explicitly define ID with correct GORM tag
	CreatedAt time.Time      // Explicitly define CreatedAt
	UpdatedAt time.Time      // Explicitly define UpdatedAt
	DeletedAt gorm.DeletedAt `gorm:"index"` // Explicitly define DeletedAt

	StaffID *uint  `gorm:"null;index"`
	Staff   *Staff `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UserID *uint `gorm:"null;index"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Action    string  `gorm:"type:varchar(255);not null"`
	Details   string  `gorm:"type:jsonb;null"`
	IPAddress *string `gorm:"type:varchar(45);null"`
	// Timestamp field from gorm.Model is CreatedAt. If a separate 'timestamp' for the log event itself is needed,
	// it should be added explicitly, e.g., EventTimestamp time.Time `gorm:"not null"`
}
