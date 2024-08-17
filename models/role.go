package models

import "github.com/google/uuid"

// User struct defines the User model
type Role struct {
	ID             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Label          string        `json:"label" binding:"required"`
	Name           string        `json:"name" binding:"required"`
	OrganizationID uuid.UUID     `gorm:"type:uuid"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID"`
	TimeMixin
}
