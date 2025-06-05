package storages

import (
	"net/http"
	"postgresus-backend/internal/features/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StorageController struct {
	storageService *StorageService
	userService    *users.UserService
}

func (c *StorageController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/storages", c.SaveStorage)
	router.GET("/storages", c.GetStorages)
	router.GET("/storages/:id", c.GetStorage)
	router.DELETE("/storages/:id", c.DeleteStorage)
	router.POST("/storages/:id/test", c.TestStorageConnection)
	router.POST("/storages/direct-test", c.TestStorageConnectionDirect)
}

// SaveStorage
// @Summary Save a storage
// @Description Create or update a storage
// @Tags storages
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param storage body Storage true "Storage data"
// @Success 200 {object} Storage
// @Failure 400
// @Failure 401
// @Router /storages [post]
func (c *StorageController) SaveStorage(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var storage Storage
	if err := ctx.ShouldBindJSON(&storage); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.storageService.SaveStorage(user, &storage); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storage)
}

// GetStorage
// @Summary Get a storage by ID
// @Description Get a specific storage by ID
// @Tags storages
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Storage ID"
// @Success 200 {object} Storage
// @Failure 400
// @Failure 401
// @Router /storages/{id} [get]
func (c *StorageController) GetStorage(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid storage ID"})
		return
	}

	storage, err := c.storageService.GetStorage(user, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storage)
}

// GetStorages
// @Summary Get all storages
// @Description Get all storages for the current user
// @Tags storages
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {array} Storage
// @Failure 401
// @Router /storages [get]
func (c *StorageController) GetStorages(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	storages, err := c.storageService.GetStorages(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, storages)
}

// DeleteStorage
// @Summary Delete a storage
// @Description Delete a storage by ID
// @Tags storages
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Storage ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /storages/{id} [delete]
func (c *StorageController) DeleteStorage(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid storage ID"})
		return
	}

	if err := c.storageService.DeleteStorage(user, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "storage deleted successfully"})
}

// TestStorageConnection
// @Summary Test storage connection
// @Description Test the connection to the storage
// @Tags storages
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Storage ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /storages/{id}/test [post]
func (c *StorageController) TestStorageConnection(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid storage ID"})
		return
	}

	if err := c.storageService.TestStorageConnection(user, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "storage connection test successful"})
}

// TestStorageConnectionDirect
// @Summary Test storage connection directly
// @Description Test the connection to a storage object provided in the request
// @Tags storages
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param storage body Storage true "Storage data"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /storages/direct-test [post]
func (c *StorageController) TestStorageConnectionDirect(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var storage Storage
	if err := ctx.ShouldBindJSON(&storage); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For direct test, associate with the current user
	storage.UserID = user.ID

	if err := storage.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.storageService.TestStorageConnectionDirect(&storage); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "storage connection test successful"})
}
