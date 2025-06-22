package databases

import (
	"errors"
	"log/slog"
	"postgresus-backend/internal/features/databases/databases/postgresql"
	"postgresus-backend/internal/features/intervals"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Database struct {
	ID          uuid.UUID    `json:"id"          gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID    `json:"userId"      gorm:"column:user_id;type:uuid;not null"`
	Name        string       `json:"name"        gorm:"column:name;type:text;not null"`
	Type        DatabaseType `json:"type"        gorm:"column:type;type:text;not null"`
	StorePeriod Period       `json:"storePeriod" gorm:"column:store_period;type:text;not null"`

	BackupIntervalID uuid.UUID           `json:"backupIntervalId"         gorm:"column:backup_interval_id;type:uuid;not null"`
	BackupInterval   *intervals.Interval `json:"backupInterval,omitempty" gorm:"foreignKey:BackupIntervalID"`

	Postgresql *postgresql.PostgresqlDatabase `json:"postgresql,omitempty" gorm:"foreignKey:DatabaseID"`

	Storage   storages.Storage `json:"storage"   gorm:"foreignKey:StorageID"`
	StorageID uuid.UUID        `json:"storageId" gorm:"column:storage_id;type:uuid;not null"`

	Notifiers                 []notifiers.Notifier     `json:"notifiers"           gorm:"many2many:database_notifiers;"`
	SendNotificationsOn       []BackupNotificationType `json:"sendNotificationsOn" gorm:"-"`
	SendNotificationsOnString string                   `json:"-"                   gorm:"column:send_notifications_on;type:text;not null"`

	// these fields are not reliable, but
	// they are used for pretty UI
	LastBackupTime         *time.Time `json:"lastBackupTime,omitempty"         gorm:"column:last_backup_time;type:timestamp with time zone"`
	LastBackupErrorMessage *string    `json:"lastBackupErrorMessage,omitempty" gorm:"column:last_backup_error_message;type:text"`
}

func (d *Database) BeforeSave(tx *gorm.DB) error {
	// Convert SendNotificationsOn array to string
	if len(d.SendNotificationsOn) > 0 {
		notificationTypes := make([]string, len(d.SendNotificationsOn))
		for i, notificationType := range d.SendNotificationsOn {
			notificationTypes[i] = string(notificationType)
		}
		d.SendNotificationsOnString = strings.Join(notificationTypes, ",")
	} else {
		d.SendNotificationsOnString = ""
	}

	return nil
}

func (d *Database) AfterFind(tx *gorm.DB) error {
	// Convert SendNotificationsOnString to array
	if d.SendNotificationsOnString != "" {
		notificationTypes := strings.Split(d.SendNotificationsOnString, ",")
		d.SendNotificationsOn = make([]BackupNotificationType, len(notificationTypes))
		for i, notificationType := range notificationTypes {
			d.SendNotificationsOn[i] = BackupNotificationType(notificationType)
		}
	} else {
		d.SendNotificationsOn = []BackupNotificationType{}
	}

	return nil
}

func (d *Database) Validate() error {
	if d.Name == "" {
		return errors.New("name is required")
	}

	// Backup interval is required either as ID or as object
	if d.BackupIntervalID == uuid.Nil && d.BackupInterval == nil {
		return errors.New("backup interval is required")
	}

	if d.StorePeriod == "" {
		return errors.New("store period is required")
	}

	if d.Postgresql.CpuCount == 0 {
		return errors.New("cpu count is required")
	}

	switch d.Type {
	case DatabaseTypePostgres:
		return d.Postgresql.Validate()
	default:
		return errors.New("invalid database type: " + string(d.Type))
	}
}

func (d *Database) ValidateUpdate(old, new Database) error {
	if old.Type != new.Type {
		return errors.New("database type is not allowed to change")
	}

	return nil
}

func (d *Database) TestConnection(logger *slog.Logger) error {
	return d.getSpecificDatabase().TestConnection(logger)
}

func (d *Database) getSpecificDatabase() DatabaseConnector {
	switch d.Type {
	case DatabaseTypePostgres:
		return d.Postgresql
	}

	panic("invalid database type: " + string(d.Type))
}
