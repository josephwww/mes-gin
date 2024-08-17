package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization struct defines the Organization model
type Organization struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"unique" json:"name" binding:"required,min=1,max=255"`
	CreatorID       uuid.UUID `gorm:"type:uuid"`                                                 // Foreign key to the Creator (User)
	Creator         *User     `gorm:"foreignKey:CreatorID"`                                      // Creator relationship
	AdministratorID uuid.UUID `json:"administrator_id" gorm:"type:uuid" binding:"required,uuid"` // Foreign key to the Administrator (User)
	Administrator   *User     `gorm:"foreignKey:AdministratorID"`                                // Administrator relationship
	TimeMixin                 // Embedding TimeMixin for common time fields
}

// CreateOrganization creates a new organization
func CreateOrganization(db *gorm.DB, org *Organization) error {
	return db.Create(org).Error
}

// GetOrganization retrieves an organization by ID
func GetOrganization(db *gorm.DB, id uuid.UUID) (*Organization, error) {
	var org Organization
	err := db.Preload("Creator").Preload("Administrator").First(&org, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// UpdateOrganization updates an existing organization
func UpdateOrganization(db *gorm.DB, id uuid.UUID, updatedOrg *Organization) error {
	var org Organization
	err := db.First(&org, "id = ?", id).Error
	if err != nil {
		return err
	}

	return db.Model(&org).Updates(updatedOrg).Error
}

// DeleteOrganization deletes an organization by ID
func DeleteOrganization(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&Organization{}, "id = ?", id).Error
}

// GetOrganizations retrieves all organizations
func GetOrganizations(db *gorm.DB) ([]Organization, error) {
	var orgs []Organization
	err := db.Preload("Creator").Preload("Administrator").Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
