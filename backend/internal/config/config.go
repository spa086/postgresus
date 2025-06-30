package config

import (
	"os"
	"path/filepath"
	env_utils "postgresus-backend/internal/util/env"
	"postgresus-backend/internal/util/logger"
	"postgresus-backend/internal/util/tools"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var log = logger.GetLogger()

const (
	AppModeWeb        = "web"
	AppModeBackground = "background"
)

type EnvVariables struct {
	IsTesting            bool
	DatabaseDsn          string            `env:"DATABASE_DSN"         required:"true"`
	EnvMode              env_utils.EnvMode `env:"ENV_MODE"             required:"true"`
	PostgresesInstallDir string            `env:"POSTGRES_INSTALL_DIR"`

	DataFolder string
	TempFolder string

	TestGoogleDriveClientID     string `env:"TEST_GOOGLE_DRIVE_CLIENT_ID"     required:"true"`
	TestGoogleDriveClientSecret string `env:"TEST_GOOGLE_DRIVE_CLIENT_SECRET" required:"true"`
	TestGoogleDriveTokenJSON    string `env:"TEST_GOOGLE_DRIVE_TOKEN_JSON"    required:"true"`
}

var (
	env  EnvVariables
	once sync.Once
)

func GetEnv() EnvVariables {
	once.Do(loadEnvVariables)
	return env
}

func loadEnvVariables() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Warn("could not get current working directory", "error", err)
		cwd = "."
	}

	backendRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(backendRoot, "go.mod")); err == nil {
			break
		}

		parent := filepath.Dir(backendRoot)
		if parent == backendRoot {
			break
		}

		backendRoot = parent
	}

	envPaths := []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(backendRoot, ".env"),
	}

	var loaded bool
	for _, path := range envPaths {
		log.Info("Trying to load .env", "path", path)
		if err := godotenv.Load(path); err == nil {
			log.Info("Successfully loaded .env", "path", path)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Error("Error loading .env file: could not find .env in any location")
		os.Exit(1)
	}

	err = cleanenv.ReadEnv(&env)
	if err != nil {
		log.Error("Configuration could not be loaded", "error", err)
		os.Exit(1)
	}

	for _, arg := range os.Args {
		if strings.Contains(arg, "test") {
			env.IsTesting = true
			break
		}
	}

	if env.DatabaseDsn == "" {
		log.Error("DATABASE_DSN is empty")
		os.Exit(1)
	}

	if env.EnvMode == "" {
		log.Error("ENV_MODE is empty")
		os.Exit(1)
	}
	if env.EnvMode != "development" && env.EnvMode != "production" {
		log.Error("ENV_MODE is invalid", "mode", env.EnvMode)
		os.Exit(1)
	}
	log.Info("ENV_MODE loaded", "mode", env.EnvMode)

	env.PostgresesInstallDir = filepath.Join(backendRoot, "tools", "postgresql")
	tools.VerifyPostgresesInstallation(log, env.EnvMode, env.PostgresesInstallDir)

	// Store the data and temp folders one level below the root
	// (projectRoot/postgresus-data -> /postgresus-data)
	env.DataFolder = filepath.Join(filepath.Dir(backendRoot), "postgresus-data", "data")
	env.TempFolder = filepath.Join(filepath.Dir(backendRoot), "postgresus-data", "temp")

	log.Info("Environment variables loaded successfully!")
}
