package config

const (
	StorageInMemory string = "in-memory"
	StorageSQL      string = "SQL"
)

var AllowedSQLDrivers [1]string = [1]string{"pgx"}

const (
	ErrFailedOpenConfigFile  = "failed to open configFil"
	ErrFailedReadConfig      = "failed to read config"
	ErrMissingDSN            = "missing DSN"
	ErrWrongSQLDriver        = "wrong SQL driver"
	ErrMissingOutputPaths    = "missing logger output paths"
	ErrMissingErrOutputPaths = "missing logger error output paths"
)
