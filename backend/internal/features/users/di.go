package users

import (
	user_repositories "postgresus-backend/internal/features/users/repositories"
)

var secretKeyRepository = &user_repositories.SecretKeyRepository{}
var userRepository = &user_repositories.UserRepository{}
var userService = &UserService{
	userRepository,
	secretKeyRepository,
}
var userController = &UserController{
	userService,
}

func GetUserService() *UserService {
	return userService
}

func GetUserController() *UserController {
	return userController
}
