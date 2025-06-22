package notifiers

import (
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
)

var notifierRepository = &NotifierRepository{}
var notifierService = &NotifierService{
	notifierRepository,
	logger.GetLogger(),
}
var notifierController = &NotifierController{
	notifierService,
	users.GetUserService(),
}

func GetNotifierController() *NotifierController {
	return notifierController
}

func GetNotifierService() *NotifierService {
	return notifierService
}
