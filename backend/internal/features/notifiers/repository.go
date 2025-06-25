package notifiers

import (
	"postgresus-backend/internal/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotifierRepository struct{}

func (r *NotifierRepository) Save(notifier *Notifier) (*Notifier, error) {
	db := storage.GetDb()

	err := db.Transaction(func(tx *gorm.DB) error {
		switch notifier.NotifierType {
		case NotifierTypeTelegram:
			if notifier.TelegramNotifier != nil {
				notifier.TelegramNotifier.NotifierID = notifier.ID
			}
		case NotifierTypeEmail:
			if notifier.EmailNotifier != nil {
				notifier.EmailNotifier.NotifierID = notifier.ID
			}
		case NotifierTypeWebhook:
			if notifier.WebhookNotifier != nil {
				notifier.WebhookNotifier.NotifierID = notifier.ID
			}
		case NotifierTypeSlack:
			if notifier.SlackNotifier != nil {
				notifier.SlackNotifier.NotifierID = notifier.ID
			}
		}

		if notifier.ID == uuid.Nil {
			if err := tx.Create(notifier).
				Omit("TelegramNotifier", "EmailNotifier", "WebhookNotifier", "SlackNotifier").
				Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(notifier).
				Omit("TelegramNotifier", "EmailNotifier", "WebhookNotifier", "SlackNotifier").
				Error; err != nil {
				return err
			}
		}

		switch notifier.NotifierType {
		case NotifierTypeTelegram:
			if notifier.TelegramNotifier != nil {
				notifier.TelegramNotifier.NotifierID = notifier.ID // Ensure ID is set
				if err := tx.Save(notifier.TelegramNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeEmail:
			if notifier.EmailNotifier != nil {
				notifier.EmailNotifier.NotifierID = notifier.ID // Ensure ID is set
				if err := tx.Save(notifier.EmailNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeWebhook:
			if notifier.WebhookNotifier != nil {
				notifier.WebhookNotifier.NotifierID = notifier.ID // Ensure ID is set
				if err := tx.Save(notifier.WebhookNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeSlack:
			if notifier.SlackNotifier != nil {
				notifier.SlackNotifier.NotifierID = notifier.ID // Ensure ID is set
				if err := tx.Save(notifier.SlackNotifier).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return notifier, nil
}

func (r *NotifierRepository) FindByID(id uuid.UUID) (*Notifier, error) {
	var notifier Notifier

	if err := storage.
		GetDb().
		Preload("TelegramNotifier").
		Preload("EmailNotifier").
		Preload("WebhookNotifier").
		Preload("SlackNotifier").
		Where("id = ?", id).
		First(&notifier).Error; err != nil {
		return nil, err
	}

	return &notifier, nil
}

func (r *NotifierRepository) FindByUserID(userID uuid.UUID) ([]*Notifier, error) {
	var notifiers []*Notifier

	if err := storage.
		GetDb().
		Preload("TelegramNotifier").
		Preload("EmailNotifier").
		Preload("WebhookNotifier").
		Preload("SlackNotifier").
		Where("user_id = ?", userID).
		Find(&notifiers).Error; err != nil {
		return nil, err
	}

	return notifiers, nil
}

func (r *NotifierRepository) Delete(notifier *Notifier) error {
	return storage.GetDb().Transaction(func(tx *gorm.DB) error {
		// Delete specific notifier based on type
		switch notifier.NotifierType {
		case NotifierTypeTelegram:
			if notifier.TelegramNotifier != nil {
				if err := tx.Delete(notifier.TelegramNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeEmail:
			if notifier.EmailNotifier != nil {
				if err := tx.Delete(notifier.EmailNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeWebhook:
			if notifier.WebhookNotifier != nil {
				if err := tx.Delete(notifier.WebhookNotifier).Error; err != nil {
					return err
				}
			}
		case NotifierTypeSlack:
			if notifier.SlackNotifier != nil {
				if err := tx.Delete(notifier.SlackNotifier).Error; err != nil {
					return err
				}
			}
		}

		// Delete the main notifier
		return tx.Delete(notifier).Error
	})
}
