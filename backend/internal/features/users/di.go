package users

import (
	user_repositories "postgresus-backend/internal/features/users/repositories"

	"golang.org/x/time/rate"
)

var secretKeyRepository = &user_repositories.SecretKeyRepository{}
var userRepository = &user_repositories.UserRepository{}
var userService = &UserService{
	userRepository,
	secretKeyRepository,
}
var userController = &UserController{
	userService,
	rate.NewLimiter(rate.Limit(3), 3), // 3 RPS with burst of 3
}

func GetUserService() *UserService {
	return userService
}

func GetUserController() *UserController {
	return userController
}
