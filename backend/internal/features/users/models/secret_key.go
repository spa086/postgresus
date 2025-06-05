package users_models

type SecretKey struct {
	Secret string `gorm:"column:secret;uniqueIndex;not null"`
}

func (SecretKey) TableName() string {
	return "secret_keys"
}
