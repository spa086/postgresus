package notifiers

import (
	"errors"
	"log/slog"
	"postgresus-backend/internal/features/notifiers/notifiers/email_notifier"
	telegram_notifier "postgresus-backend/internal/features/notifiers/notifiers/telegram"
	webhook_notifier "postgresus-backend/internal/features/notifiers/notifiers/webhook"

	"github.com/google/uuid"
)

type Notifier struct {
	ID            uuid.UUID    `json:"id"            gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        uuid.UUID    `json:"userId"        gorm:"column:user_id;not null;type:uuid;index"`
	Name          string       `json:"name"          gorm:"column:name;not null;type:varchar(255)"`
	NotifierType  NotifierType `json:"notifierType"  gorm:"column:notifier_type;not null;type:varchar(50)"`
	LastSendError *string      `json:"lastSendError" gorm:"column:last_send_error;type:text"`

	// specific notifier
	TelegramNotifier *telegram_notifier.TelegramNotifier `json:"telegramNotifier" gorm:"foreignKey:NotifierID"`
	EmailNotifier    *email_notifier.EmailNotifier       `json:"emailNotifier"    gorm:"foreignKey:NotifierID"`
	WebhookNotifier  *webhook_notifier.WebhookNotifier   `json:"webhookNotifier"  gorm:"foreignKey:NotifierID"`
}

func (n *Notifier) TableName() string {
	return "notifiers"
}

func (n *Notifier) Validate() error {
	if n.Name == "" {
		return errors.New("name is required")
	}

	return n.getSpecificNotifier().Validate()
}

func (n *Notifier) Send(logger *slog.Logger, heading string, message string) error {
	err := n.getSpecificNotifier().Send(logger, heading, message)

	if err != nil {
		lastSendError := err.Error()
		n.LastSendError = &lastSendError
	} else {
		n.LastSendError = nil
	}

	return err
}

func (n *Notifier) getSpecificNotifier() NotificationSender {
	switch n.NotifierType {
	case NotifierTypeTelegram:
		return n.TelegramNotifier
	case NotifierTypeEmail:
		return n.EmailNotifier
	case NotifierTypeWebhook:
		return n.WebhookNotifier
	default:
		panic("unknown notifier type: " + string(n.NotifierType))
	}
}
