package databases

import (
	"errors"
	users_models "postgresus-backend/internal/features/users/models"
	"time"

	"github.com/google/uuid"
)

type DatabaseService struct {
	dbRepository            *DatabaseRepository
	dbStorageChangeListener DatabaseStorageChangeListener
}

func (s *DatabaseService) SetDatabaseStorageChangeListener(
	dbStorageChangeListener DatabaseStorageChangeListener,
) {
	s.dbStorageChangeListener = dbStorageChangeListener
}

func (s *DatabaseService) CreateDatabase(
	user *users_models.User,
	database *Database,
) error {
	database.UserID = user.ID

	if err := database.Validate(); err != nil {
		return err
	}

	return s.dbRepository.Save(database)
}

func (s *DatabaseService) UpdateDatabase(
	user *users_models.User,
	database *Database,
) error {
	if database.ID == uuid.Nil {
		return errors.New("database ID is required for update")
	}

	existingDatabase, err := s.dbRepository.FindByID(database.ID)
	if err != nil {
		return err
	}

	if existingDatabase.UserID != user.ID {
		return errors.New("you have not access to this database")
	}

	// Validate the update
	if err := database.ValidateUpdate(*existingDatabase, *database); err != nil {
		return err
	}

	if err := database.Validate(); err != nil {
		return err
	}

	if existingDatabase.Storage.ID != database.Storage.ID {
		err := s.dbStorageChangeListener.OnBeforeDbStorageChange(
			existingDatabase.ID,
			database.StorageID,
		)

		if err != nil {
			return err
		}
	}

	return s.dbRepository.Save(database)
}

func (s *DatabaseService) DeleteDatabase(
	user *users_models.User,
	id uuid.UUID,
) error {
	existingDatabase, err := s.dbRepository.FindByID(id)
	if err != nil {
		return err
	}

	if existingDatabase.UserID != user.ID {
		return errors.New("you have not access to this database")
	}

	return s.dbRepository.Delete(id)
}

func (s *DatabaseService) GetDatabase(
	user *users_models.User,
	id uuid.UUID,
) (*Database, error) {
	database, err := s.dbRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if database.UserID != user.ID {
		return nil, errors.New("you have not access to this database")
	}

	return database, nil
}

func (s *DatabaseService) GetDatabasesByUser(
	user *users_models.User,
) ([]*Database, error) {
	return s.dbRepository.FindByUserID(user.ID)
}

func (s *DatabaseService) IsNotifierUsing(notifierID uuid.UUID) (bool, error) {
	return s.dbRepository.IsNotifierUsing(notifierID)
}

func (s *DatabaseService) IsStorageUsing(storageID uuid.UUID) (bool, error) {
	return s.dbRepository.IsStorageUsing(storageID)
}

func (s *DatabaseService) TestDatabaseConnection(
	user *users_models.User,
	databaseID uuid.UUID,
) error {
	database, err := s.dbRepository.FindByID(databaseID)
	if err != nil {
		return err
	}

	if database.UserID != user.ID {
		return errors.New("you have not access to this database")
	}

	err = database.TestConnection()
	if err != nil {
		lastSaveError := err.Error()
		database.LastBackupErrorMessage = &lastSaveError
		return err
	}

	database.LastBackupErrorMessage = nil

	return s.dbRepository.Save(database)
}

func (s *DatabaseService) TestDatabaseConnectionDirect(
	database *Database,
) error {
	return database.TestConnection()
}

func (s *DatabaseService) GetDatabaseByID(
	id uuid.UUID,
) (*Database, error) {
	return s.dbRepository.FindByID(id)
}

func (s *DatabaseService) GetAllDatabases() ([]*Database, error) {
	return s.dbRepository.GetAllDatabases()
}

func (s *DatabaseService) SetBackupError(databaseID uuid.UUID, errorMessage string) error {
	database, err := s.dbRepository.FindByID(databaseID)
	if err != nil {
		return err
	}

	database.LastBackupErrorMessage = &errorMessage
	return s.dbRepository.Save(database)
}

func (s *DatabaseService) SetLastBackupTime(databaseID uuid.UUID, backupTime time.Time) error {
	database, err := s.dbRepository.FindByID(databaseID)
	if err != nil {
		return err
	}

	database.LastBackupTime = &backupTime
	database.LastBackupErrorMessage = nil // Clear any previous error
	return s.dbRepository.Save(database)
}
