package handler

import (
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mes-gin/config"
	"mes-gin/models"
)

// GetOrganizations handles retrieving all organizations
func GetOrganizations(c *gin.Context) {
	orgs, err := models.GetOrganizations(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organizations"})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// GetOrganization handles the retrieval of an organization by ID
func GetOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := models.GetOrganization(config.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// CreateOrganization handles the creation of a new organization
func CreateOrganization(c *gin.Context) {
	// Retrieve currentUser from the context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Cast currentUser to the User model
	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	var org models.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if Administrator User exists
	if _, err := models.GetUserByID(config.DB, org.AdministratorID.String()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org.CreatorID = user.ID

	if err := models.CreateOrganization(config.DB, &org); err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok && pqErr.Code == "23505" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": pqErr.Detail})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
			return
		}
	}

	c.JSON(http.StatusCreated, org)
}

// UpdateOrganization handles the update of an existing organization
func UpdateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var updatedOrg models.Organization
	if err := c.ShouldBindJSON(&updatedOrg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateOrganization(config.DB, id, &updatedOrg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, updatedOrg)
}

// DeleteOrganization handles the deletion of an organization by ID
func DeleteOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := models.DeleteOrganization(config.DB, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted"})
}
