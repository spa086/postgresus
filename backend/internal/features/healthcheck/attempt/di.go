package healthcheck_attempt

import (
	"postgresus-backend/internal/features/databases"
	healthcheck_config "postgresus-backend/internal/features/healthcheck/config"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
)

var healthcheckAttemptRepository = &HealthcheckAttemptRepository{}
var healthcheckAttemptService = &HealthcheckAttemptService{
	healthcheckAttemptRepository,
	databases.GetDatabaseService(),
}

var checkPgHealthUseCase = &CheckPgHealthUseCase{
	healthcheckAttemptRepository,
	notifiers.GetNotifierService(),
	databases.GetDatabaseService(),
}

var healthcheckAttemptBackgroundService = &HealthcheckAttemptBackgroundService{
	healthcheck_config.GetHealthcheckConfigService(),
	checkPgHealthUseCase,
	logger.GetLogger(),
}
var healthcheckAttemptController = &HealthcheckAttemptController{
	healthcheckAttemptService,
	users.GetUserService(),
}

func GetHealthcheckAttemptService() *HealthcheckAttemptService {
	return healthcheckAttemptService
}

func GetHealthcheckAttemptBackgroundService() *HealthcheckAttemptBackgroundService {
	return healthcheckAttemptBackgroundService
}

func GetHealthcheckAttemptController() *HealthcheckAttemptController {
	return healthcheckAttemptController
}
