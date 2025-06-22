package notifiers

import (
	"errors"
	"log/slog"
	users_models "postgresus-backend/internal/features/users/models"

	"github.com/google/uuid"
)

type NotifierService struct {
	notifierRepository *NotifierRepository
	logger             *slog.Logger
}

func (s *NotifierService) SaveNotifier(
	user *users_models.User,
	notifier *Notifier,
) error {
	if notifier.ID != uuid.Nil {
		existingNotifier, err := s.notifierRepository.FindByID(notifier.ID)
		if err != nil {
			return err
		}

		if existingNotifier.UserID != user.ID {
			return errors.New("you have not access to this notifier")
		}

		notifier.UserID = existingNotifier.UserID
	} else {
		notifier.UserID = user.ID
	}

	return s.notifierRepository.Save(notifier)
}

func (s *NotifierService) DeleteNotifier(
	user *users_models.User,
	notifierID uuid.UUID,
) error {
	notifier, err := s.notifierRepository.FindByID(notifierID)
	if err != nil {
		return err
	}

	if notifier.UserID != user.ID {
		return errors.New("you have not access to this notifier")
	}

	return s.notifierRepository.Delete(notifier)
}

func (s *NotifierService) GetNotifier(
	user *users_models.User,
	id uuid.UUID,
) (*Notifier, error) {
	notifier, err := s.notifierRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if notifier.UserID != user.ID {
		return nil, errors.New("you have not access to this notifier")
	}

	return notifier, nil
}

func (s *NotifierService) GetNotifiers(
	user *users_models.User,
) ([]*Notifier, error) {
	return s.notifierRepository.FindByUserID(user.ID)
}

func (s *NotifierService) SendTestNotification(
	user *users_models.User,
	notifierID uuid.UUID,
) error {
	notifier, err := s.notifierRepository.FindByID(notifierID)
	if err != nil {
		return err
	}

	if notifier.UserID != user.ID {
		return errors.New("you have not access to this notifier")
	}

	err = notifier.Send(s.logger, "Test message", "This is a test message")
	if err != nil {
		return err
	}

	if err = s.notifierRepository.Save(notifier); err != nil {
		return err
	}

	return nil
}

func (s *NotifierService) SendTestNotificationToNotifier(
	notifier *Notifier,
) error {
	return notifier.Send(s.logger, "Test message", "This is a test message")
}

func (s *NotifierService) SendNotification(
	notifier *Notifier,
	title string,
	message string,
) {
	// Truncate message to 2000 characters if it's too long
	messageRunes := []rune(message)
	if len(messageRunes) > 2000 {
		message = string(messageRunes[:2000])
	}

	notifiedFromDb, err := s.notifierRepository.FindByID(notifier.ID)
	if err != nil {
		return
	}

	err = notifiedFromDb.Send(s.logger, title, message)
	if err != nil {
		errMsg := err.Error()
		notifiedFromDb.LastSendError = &errMsg

		err = s.notifierRepository.Save(notifiedFromDb)
		if err != nil {
			s.logger.Error("Failed to save notifier", "error", err)
		}
	}

	notifiedFromDb.LastSendError = nil
	err = s.notifierRepository.Save(notifiedFromDb)
	if err != nil {
		s.logger.Error("Failed to save notifier", "error", err)
	}
}
