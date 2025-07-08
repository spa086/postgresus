package databases

type DatabaseType string

const (
	DatabaseTypePostgres DatabaseType = "POSTGRES"
)

type HealthStatus string

const (
	HealthStatusAvailable   HealthStatus = "AVAILABLE"
	HealthStatusUnavailable HealthStatus = "UNAVAILABLE"
)
