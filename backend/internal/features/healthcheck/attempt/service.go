package healthcheck_attempt

import (
	"errors"
	"postgresus-backend/internal/features/databases"
	users_models "postgresus-backend/internal/features/users/models"
	"time"

	"github.com/google/uuid"
)

type HealthcheckAttemptService struct {
	healthcheckAttemptRepository *HealthcheckAttemptRepository
	databaseService              *databases.DatabaseService
}

func (s *HealthcheckAttemptService) GetAttemptsByDatabase(
	user users_models.User,
	databaseID uuid.UUID,
	afterDate time.Time,
) ([]*HealthcheckAttempt, error) {
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return nil, err
	}

	if database.UserID != user.ID {
		return nil, errors.New("forbidden")
	}

	return s.healthcheckAttemptRepository.FindByDatabaseIdOrderByCreatedAtDesc(
		databaseID,
		afterDate,
	)
}
