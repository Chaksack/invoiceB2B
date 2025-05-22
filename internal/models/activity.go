package models

import "gorm.io/gorm"

//type ActivityLog struct {
//	ID        uint `gorm:"primarykey"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt gorm.DeletedAt `gorm:"index"`
//
//	StaffID *uint  `gorm:"null;index"`
//	Staff   *Staff `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//
//	UserID *uint `gorm:"null;index"`
//	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//
//	Action    string  `gorm:"type:varchar(255);not null"`
//	Details   string  `gorm:"type:jsonb;null"`
//	IPAddress *string `gorm:"type:varchar(45);null"`
//}

type ActivityLog struct {
	gorm.Model

	StaffID *uint  `gorm:"null;index"`
	Staff   *Staff `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UserID *uint `gorm:"null;index"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Action    string  `gorm:"type:varchar(255);not null"`
	Details   string  `gorm:"type:jsonb;null"` // Can store marshalled JSON like {"invoice_id": 123, "old_status": "pending", "new_status": "approved"}
	IPAddress *string `gorm:"type:varchar(45);null"`
}
