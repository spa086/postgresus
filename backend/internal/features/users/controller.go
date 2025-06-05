package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *UserService
}

func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/users/signup", c.SignUp)
	router.POST("/users/signin", c.SignIn)
	router.GET("/users/is-any-user-exist", c.IsAnyUserExist)
}

// SignUp
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "User signup data"
// @Success 200
// @Failure 400
// @Router /users/signup [post]
func (c *UserController) SignUp(ctx *gin.Context) {
	var request SignUpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := c.userService.SignUp(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// SignIn
// @Summary Authenticate a user
// @Description Authenticate a user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignInRequest true "User signin data"
// @Success 200 {object} SignInResponse
// @Failure 400
// @Router /users/signin [post]
func (c *UserController) SignIn(ctx *gin.Context) {
	var request SignInRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := c.userService.SignIn(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// IsAnyUserExist
// @Summary Check if any user exists
// @Description Check if any user exists in the system
// @Tags users
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /users/is-any-user-exist [get]
func (c *UserController) IsAnyUserExist(ctx *gin.Context) {
	isExist, err := c.userService.IsAnyUserExist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"isExist": isExist})
}
