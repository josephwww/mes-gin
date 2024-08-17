package models

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"mes-gin/api/v1/request"
)

// User struct defines the User model
type User struct {
	ID             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string        `json:"name" binding:"required,min=1,max=100"` // Name is required and must be between 3 and 100 characters
	Phone          string        `json:"phone" binding:"required"`              // Email is required and must be a valid email format
	Password       string        `json:"password" binding:"required,min=8"`     // Password is required and must be at least 8 characters
	OrganizationID *uuid.UUID    `gorm:"type:uuid" binding:"required"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID"`           // Organization relationship
	Roles          []*Role       `gorm:"many2many:user_roles;" json:"roles"`  // Many-to-many relationship with Role
	CreatorID      *uuid.UUID    `gorm:"type:uuid" json:"creator_id"`         // Foreign key for the user who created this user
	Creator        *User         `gorm:"foreignKey:CreatorID" json:"creator"` // Creator relationship
	TimeMixin
}
type RoleClaims struct {
	RoleID    uuid.UUID
	RoleLabel string
	RoleName  string
}

// Claims defines the structure for JWT claims
type Claims struct {
	UserID         uuid.UUID    `json:"user_id"`
	UserName       string       `json:"username"`
	OrganizationID *uuid.UUID   `json:"organization_id"`
	Roles          []RoleClaims `json:"roles"`
	jwt.StandardClaims
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByID fetches a user by their ID
func GetUserByID(db *gorm.DB, userID string) (*User, error) {
	var user User
	result := db.First(&user, "id = ?", userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New(fmt.Sprintf("user_id: %v not found", userID))
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetAllUsers fetches all users from the database
func GetAllUsers(db *gorm.DB, listParams request.Pagination) ([]User, int64, error) {
	var users []User
	var users_count int64

	base_query := db
	if &listParams.Query != nil {
		base_query = base_query.Where(
			"name LIKE ? or phone LIKE ?",
			fmt.Sprintf("%%%v%%", listParams.Query),
			fmt.Sprintf("%%%v%%", listParams.Query),
		)
	}

	result := base_query.Find(&users).Count(&users_count)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = base_query.Offset(*listParams.Start).Limit(*listParams.Limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, users_count, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(db *gorm.DB, user *User) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
