package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// ActivityLog records key actions within the system
type ActivityLog struct {
	ID      uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	StaffID *uuid.UUID `gorm:"type:uuid;null;index"` // FK to Staff, nullable for system actions
	Staff   *Staff     `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UserID *uuid.UUID `gorm:"type:uuid;null;index"` // FK to User, nullable for non-user specific actions
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Action    string    `gorm:"type:varchar(255);not null"` // e.g., 'USER_REGISTERED', 'KYC_APPROVED'
	Details   string    `gorm:"type:jsonb;null"`            // Additional context as JSON string
	IPAddress *string   `gorm:"type:varchar(45);null"`
	Timestamp time.Time `gorm:"autoCreateTime;not null"`
}

// BeforeCreate hook for ActivityLog to generate UUID
func (al *ActivityLog) BeforeCreate(tx *gorm.DB) (err error) {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	if al.Timestamp.IsZero() {
		al.Timestamp = time.Now()
	}
	return
}
