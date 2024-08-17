package models

import (
	"gorm.io/gorm"
	"time"
)

// TimeMixin provides common time fields for models
type TimeMixin struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
