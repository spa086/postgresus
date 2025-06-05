package notifiers

import (
	"net/http"
	"postgresus-backend/internal/features/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotifierController struct {
	notifierService *NotifierService
	userService     *users.UserService
}

func (c *NotifierController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/notifiers", c.SaveNotifier)
	router.GET("/notifiers", c.GetNotifiers)
	router.GET("/notifiers/:id", c.GetNotifier)
	router.DELETE("/notifiers/:id", c.DeleteNotifier)
	router.POST("/notifiers/:id/test", c.SendTestNotification)
	router.POST("/notifiers/direct-test", c.SendTestNotificationDirect)
}

// SaveNotifier
// @Summary Save a notifier
// @Description Create or update a notifier
// @Tags notifiers
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param notifier body Notifier true "Notifier data"
// @Success 200 {object} Notifier
// @Failure 400
// @Failure 401
// @Router /notifiers [post]
func (c *NotifierController) SaveNotifier(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var notifier Notifier
	if err := ctx.ShouldBindJSON(&notifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := notifier.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.notifierService.SaveNotifier(user, &notifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifier)
}

// GetNotifier
// @Summary Get a notifier by ID
// @Description Get a specific notifier by ID
// @Tags notifiers
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Notifier ID"
// @Success 200 {object} Notifier
// @Failure 400
// @Failure 401
// @Router /notifiers/{id} [get]
func (c *NotifierController) GetNotifier(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notifier ID"})
		return
	}

	notifier, err := c.notifierService.GetNotifier(user, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifier)
}

// GetNotifiers
// @Summary Get all notifiers
// @Description Get all notifiers for the current user
// @Tags notifiers
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {array} Notifier
// @Failure 401
// @Router /notifiers [get]
func (c *NotifierController) GetNotifiers(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	notifiers, err := c.notifierService.GetNotifiers(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifiers)
}

// DeleteNotifier
// @Summary Delete a notifier
// @Description Delete a notifier by ID
// @Tags notifiers
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Notifier ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /notifiers/{id} [delete]
func (c *NotifierController) DeleteNotifier(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notifier ID"})
		return
	}

	notifier, err := c.notifierService.GetNotifier(user, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.notifierService.DeleteNotifier(user, notifier.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "notifier deleted successfully"})
}

// SendTestNotification
// @Summary Send test notification
// @Description Send a test notification using the specified notifier
// @Tags notifiers
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Notifier ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /notifiers/{id}/test [post]
func (c *NotifierController) SendTestNotification(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notifier ID"})
		return
	}

	if err := c.notifierService.SendTestNotification(user, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "test notification sent successfully"})
}

// SendTestNotificationDirect
// @Summary Send test notification directly
// @Description Send a test notification using a notifier object provided in the request
// @Tags notifiers
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param notifier body Notifier true "Notifier data"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /notifiers/direct-test [post]
func (c *NotifierController) SendTestNotificationDirect(ctx *gin.Context) {
	user, err := c.userService.GetUserFromToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var notifier Notifier
	if err := ctx.ShouldBindJSON(&notifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For direct test, associate with the current user
	notifier.UserID = user.ID

	if err := c.notifierService.SendTestNotificationToNotifier(&notifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "test notification sent successfully"})
}
