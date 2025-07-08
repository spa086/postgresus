package usecases_postgresql

import (
	"postgresus-backend/internal/util/logger"
)

var createPostgresqlBackupUsecase = &CreatePostgresqlBackupUsecase{
	logger.GetLogger(),
}

func GetCreatePostgresqlBackupUsecase() *CreatePostgresqlBackupUsecase {
	return createPostgresqlBackupUsecase
}
