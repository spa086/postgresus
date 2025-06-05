package notifiers

import "postgresus-backend/internal/features/users"

var notifierRepository = &NotifierRepository{}
var notifierService = &NotifierService{
	notifierRepository,
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
