package storages

type StorageType string

const (
	StorageTypeLocal StorageType = "LOCAL"
	StorageTypeS3    StorageType = "S3"
)
