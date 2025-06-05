package databases

import "postgresus-backend/internal/features/users"

var databaseRepository = &DatabaseRepository{}

var databaseService = &DatabaseService{
	databaseRepository,
	nil,
}

var databaseController = &DatabaseController{
	databaseService,
	users.GetUserService(),
}

func GetDatabaseService() *DatabaseService {
	return databaseService
}

func GetDatabaseController() *DatabaseController {
	return databaseController
}
