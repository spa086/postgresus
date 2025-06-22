package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"postgresus-backend/internal/config"
	"postgresus-backend/internal/features/backups"
	usecases_postgresql_backup "postgresus-backend/internal/features/backups/usecases/postgresql"
	"postgresus-backend/internal/features/databases"
	pgtypes "postgresus-backend/internal/features/databases/databases/postgresql"
	"postgresus-backend/internal/features/restores/models"
	usecases_postgresql_restore "postgresus-backend/internal/features/restores/usecases/postgresql"
	"postgresus-backend/internal/features/storages"
	local_storage "postgresus-backend/internal/features/storages/storages/local"
	"postgresus-backend/internal/util/tools"
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const createAndFillTableQuery = `
CREATE TABLE test_data (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO test_data (name, value) VALUES
    ('test1', 100),
    ('test2', 200),
    ('test3', 300);
`

type PostgresContainer struct {
	Container testcontainers.Container
	Host      string
	Port      int
	Username  string
	Password  string
	Database  string
	Version   string
	DB        *sqlx.DB
}

type TestDataItem struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Value     int       `db:"value"`
	CreatedAt time.Time `db:"created_at"`
}

// Main test functions for each PostgreSQL version
func Test_BackupAndRestorePostgresql_RestoreIsSuccesful(t *testing.T) {
	cases := []struct {
		name    string
		version string
	}{
		{"PostgreSQL 13", "13"},
		{"PostgreSQL 14", "14"},
		{"PostgreSQL 15", "15"},
		{"PostgreSQL 16", "16"},
		{"PostgreSQL 17", "17"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testBackupRestoreForVersion(t, tc.version)
		})
	}
}

// Run a test for a specific PostgreSQL version
func testBackupRestoreForVersion(t *testing.T, pgVersion string) {
	ctx := context.Background()

	// Start PostgreSQL container
	container, err := startPostgresContainer(ctx, pgVersion)
	assert.NoError(t, err)
	defer func() {
		if container.DB != nil {
			container.DB.Close()
		}

		if container.Container != nil {
			container.Container.Terminate(ctx)
		}
	}()

	_, err = container.DB.Exec(createAndFillTableQuery)
	assert.NoError(t, err)

	// Prepare data for backup
	backupID := uuid.New()
	pgVersionEnum := tools.GetPostgresqlVersionEnum(pgVersion)

	backupDbConfig := &databases.Database{
		ID:   uuid.New(),
		Type: databases.DatabaseTypePostgres,
		Name: "Test Database",
		Postgresql: &pgtypes.PostgresqlDatabase{
			Version:  pgVersionEnum,
			Host:     container.Host,
			Port:     container.Port,
			Username: container.Username,
			Password: container.Password,
			Database: &container.Database,
			IsHttps:  false,
			CpuCount: 1,
		},
	}

	storage := &storages.Storage{
		UserID:       uuid.New(),
		Type:         storages.StorageTypeLocal,
		Name:         "Test Storage",
		LocalStorage: &local_storage.LocalStorage{},
	}

	// Make backup
	progressTracker := func(completedMBs float64) {}
	err = usecases_postgresql_backup.GetCreatePostgresqlBackupUsecase().Execute(
		backupID,
		backupDbConfig,
		storage,
		progressTracker,
	)
	assert.NoError(t, err)

	// Create new database
	newDBName := "restoreddb"
	_, err = container.DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", newDBName))
	assert.NoError(t, err)

	// Connect to the new database
	newDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		container.Host, container.Port, container.Username, container.Password, newDBName)
	newDB, err := sqlx.Connect("postgres", newDSN)
	assert.NoError(t, err)
	defer newDB.Close()

	// Setup data for restore
	completedBackup := &backups.Backup{
		ID:         backupID,
		DatabaseID: backupDbConfig.ID,
		StorageID:  storage.ID,
		Status:     backups.BackupStatusCompleted,
		CreatedAt:  time.Now().UTC(),
		Storage:    storage,
		Database:   backupDbConfig,
	}

	restoreID := uuid.New()
	restore := models.Restore{
		ID:     restoreID,
		Backup: completedBackup,
		Postgresql: &pgtypes.PostgresqlDatabase{
			Version:  pgVersionEnum,
			Host:     container.Host,
			Port:     container.Port,
			Username: container.Username,
			Password: container.Password,
			Database: &newDBName,
			IsHttps:  false,
			CpuCount: 1,
		},
	}

	// Restore the backup
	restoreBackupUC := usecases_postgresql_restore.GetRestorePostgresqlBackupUsecase()
	err = restoreBackupUC.Execute(restore, completedBackup, storage)
	assert.NoError(t, err)

	// Verify restored table exists
	var tableExists bool
	err = newDB.Get(
		&tableExists,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'test_data')",
	)
	assert.NoError(t, err)
	assert.True(t, tableExists, "Table 'test_data' should exist in restored database")

	// Verify data integrity
	verifyDataIntegrity(t, container.DB, newDB)

	// Clean up the backup file after the test
	err = os.Remove(filepath.Join(config.GetEnv().DataFolder, backupID.String()))
	if err != nil {
		t.Logf("Warning: Failed to delete backup file: %v", err)
	}
}

// verifyDataIntegrity compares data in the original and restored databases
func verifyDataIntegrity(t *testing.T, originalDB *sqlx.DB, restoredDB *sqlx.DB) {
	var originalData []TestDataItem
	var restoredData []TestDataItem

	err := originalDB.Select(&originalData, "SELECT * FROM test_data ORDER BY id")
	assert.NoError(t, err)

	err = restoredDB.Select(&restoredData, "SELECT * FROM test_data ORDER BY id")
	assert.NoError(t, err)

	assert.Equal(t, len(originalData), len(restoredData), "Should have same number of rows")
	for i := range originalData {
		assert.Equal(t, originalData[i].ID, restoredData[i].ID, "ID should match")
		assert.Equal(t, originalData[i].Name, restoredData[i].Name, "Name should match")
		assert.Equal(t, originalData[i].Value, restoredData[i].Value, "Value should match")
	}
}

func startPostgresContainer(ctx context.Context, version string) (*PostgresContainer, error) {
	dbName := "testdb"
	password := "postgres"
	username := "postgres"
	port := "5432/tcp"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("postgres:%s", version),
		ExposedPorts: []string{port},
		Env: map[string]string{
			"POSTGRES_PASSWORD": password,
			"POSTGRES_USER":     username,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(port)),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedHost, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, err
	}

	portInt, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to parse port: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		mappedHost, portInt, username, password, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresContainer{
		Container: container,
		Host:      mappedHost,
		Port:      portInt,
		Username:  username,
		Password:  password,
		Database:  dbName,
		Version:   version,
		DB:        db,
	}, nil
}
