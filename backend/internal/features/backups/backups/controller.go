package backups

import (
	"fmt"
	"io"
	"net/http"
	"postgresus-backend/internal/features/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BackupController struct {
	backupService *BackupService
	userService   *users.UserService
}

func (c *BackupController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/backups", c.GetBackups)
	router.POST("/backups", c.MakeBackup)
	router.GET("/backups/:id/file", c.GetFile)
	router.DELETE("/backups/:id", c.DeleteBackup)
}

// GetBackups
// @Summary Get backups for a database
// @Description Get all backups for the specified database
// @Tags backups
// @Produce json
// @Param database_id query string true "Database ID"
// @Success 200 {array} Backup
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /backups [get]
func (c *BackupController) GetBackups(ctx *gin.Context) {
	databaseIDStr := ctx.Query("database_id")
	if databaseIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "database_id query parameter is required"})
		return
	}

	databaseID, err := uuid.Parse(databaseIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid database_id"})
		return
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	user, err := c.userService.GetUserFromToken(authorizationHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	backups, err := c.backupService.GetBackups(user, databaseID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, backups)
}

// MakeBackup
// @Summary Create a backup
// @Description Create a new backup for the specified database
// @Tags backups
// @Accept json
// @Produce json
// @Param request body MakeBackupRequest true "Backup creation data"
// @Success 200 {object} map[string]string
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /backups [post]
func (c *BackupController) MakeBackup(ctx *gin.Context) {
	var request MakeBackupRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	user, err := c.userService.GetUserFromToken(authorizationHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if err := c.backupService.MakeBackupWithAuth(user, request.DatabaseID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "backup started successfully"})
}

// DeleteBackup
// @Summary Delete a backup
// @Description Delete an existing backup
// @Tags backups
// @Param id path string true "Backup ID"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /backups/{id} [delete]
func (c *BackupController) DeleteBackup(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup ID"})
		return
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	user, err := c.userService.GetUserFromToken(authorizationHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if err := c.backupService.DeleteBackup(user, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetFile
// @Summary Download a backup file
// @Description Download the backup file for the specified backup
// @Tags backups
// @Param id path string true "Backup ID"
// @Success 200 {file} file
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /backups/{id}/file [get]
func (c *BackupController) GetFile(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup ID"})
		return
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	user, err := c.userService.GetUserFromToken(authorizationHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	fileReader, err := c.backupService.GetBackupFile(user, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := fileReader.Close(); err != nil {
			// Log the error but don't interrupt the response
			fmt.Printf("Error closing file reader: %v\n", err)
		}
	}()

	// Set headers for file download
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"backup_%s.dump\"", id.String()),
	)

	// Stream the file content
	_, err = io.Copy(ctx.Writer, fileReader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stream file"})
		return
	}
}

type MakeBackupRequest struct {
	DatabaseID uuid.UUID `json:"database_id" binding:"required"`
}
