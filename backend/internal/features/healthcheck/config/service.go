package healthcheck_config

import (
	"errors"
	"log/slog"
	"postgresus-backend/internal/features/databases"
	users_models "postgresus-backend/internal/features/users/models"

	"github.com/google/uuid"
)

type HealthcheckConfigService struct {
	databaseService             *databases.DatabaseService
	healthcheckConfigRepository *HealthcheckConfigRepository
	logger                      *slog.Logger
}

func (s *HealthcheckConfigService) OnDatabaseCreated(
	databaseID uuid.UUID,
) {
	err := s.initializeDefaultConfig(databaseID)
	if err != nil {
		s.logger.Error("failed to initialize default healthcheck config", "error", err)
	}
}

func (s *HealthcheckConfigService) Save(
	user users_models.User,
	configDTO HealthcheckConfigDTO,
) error {
	database, err := s.databaseService.GetDatabaseByID(configDTO.DatabaseID)
	if err != nil {
		return err
	}

	if database.UserID != user.ID {
		return errors.New("user does not have access to this database")
	}

	healthcheckConfig := configDTO.ToDTO()
	s.logger.Info("healthcheck config", "config", healthcheckConfig)

	healthcheckConfig.DatabaseID = database.ID

	err = s.healthcheckConfigRepository.Save(healthcheckConfig)
	if err != nil {
		return err
	}

	// for DBs with disabled healthcheck, we keep
	// health status as available
	if !healthcheckConfig.IsHealthcheckEnabled &&
		database.HealthStatus != nil {
		err = s.databaseService.SetHealthStatus(
			database.ID,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HealthcheckConfigService) GetByDatabaseID(
	user users_models.User,
	databaseID uuid.UUID,
) (*HealthcheckConfig, error) {
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return nil, err
	}

	if database.UserID != user.ID {
		return nil, errors.New("user does not have access to this database")
	}

	config, err := s.healthcheckConfigRepository.GetByDatabaseID(database.ID)
	if err != nil {
		return nil, err
	}

	if config == nil {
		err = s.initializeDefaultConfig(database.ID)
		if err != nil {
			return nil, err
		}

		config, err = s.healthcheckConfigRepository.GetByDatabaseID(database.ID)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (s *HealthcheckConfigService) GetDatabasesWithEnabledHealthcheck() (
	[]HealthcheckConfig, error,
) {
	return s.healthcheckConfigRepository.GetDatabasesWithEnabledHealthcheck()
}

func (s *HealthcheckConfigService) initializeDefaultConfig(
	databaseID uuid.UUID,
) error {
	return s.healthcheckConfigRepository.Save(&HealthcheckConfig{
		DatabaseID:                        databaseID,
		IsHealthcheckEnabled:              true,
		IsSentNotificationWhenUnavailable: true,
		IntervalMinutes:                   1,
		AttemptsBeforeConcideredAsDown:    3,
		StoreAttemptsDays:                 7,
	})
}
