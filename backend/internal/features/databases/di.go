package databases

import (
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
)

var databaseRepository = &DatabaseRepository{}

var databaseService = &DatabaseService{
	databaseRepository,
	nil,
	logger.GetLogger(),
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
