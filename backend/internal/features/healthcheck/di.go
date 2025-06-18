package healthcheck

import "postgresus-backend/internal/features/disk"

var (
	healthcheckService    *HealthcheckService
	healthcheckController *HealthcheckController
)

func init() {
	healthcheckService = &HealthcheckService{
		disk.GetDiskService(),
	}

	healthcheckController = &HealthcheckController{
		healthcheckService,
	}
}

func GetHealthcheckController() *HealthcheckController {
	return healthcheckController
}
