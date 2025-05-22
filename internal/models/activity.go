package models

import (
	"gorm.io/gorm"
	"time"
)

type ActivityLog struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	StaffID *uint  `gorm:"null;index"`
	Staff   *Staff `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UserID *uint `gorm:"null;index"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Action    string  `gorm:"type:varchar(255);not null"`
	Details   string  `gorm:"type:jsonb;null"`
	IPAddress *string `gorm:"type:varchar(45);null"`
}
