package model

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	AccountID  int            `gorm:"column:account_id"`
	LocationID int            `gorm:"column:location_id"`
	Status     string         `gorm:"column:status"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
    UpdatedAt  time.Time      `gorm:"column:updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Attendance) TableName() string {
	return "attendances"
}
