package healthcheck_config

import (
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
)

var healthcheckConfigRepository = &HealthcheckConfigRepository{}
var healthcheckConfigService = &HealthcheckConfigService{
	databases.GetDatabaseService(),
	healthcheckConfigRepository,
	logger.GetLogger(),
}
var healthcheckConfigController = &HealthcheckConfigController{
	healthcheckConfigService,
	users.GetUserService(),
}

func GetHealthcheckConfigService() *HealthcheckConfigService {
	return healthcheckConfigService
}

func GetHealthcheckConfigController() *HealthcheckConfigController {
	return healthcheckConfigController
}

func SetupDependencies() {
	databases.
		GetDatabaseService().
		AddDbCreationListener(healthcheckConfigService)
}
