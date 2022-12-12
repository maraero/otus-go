package config

const (
	StorageInMemory string = "in-memory"
	StorageSQL      string = "sql"
)

var AllowedDatabases = [1]string{"postgres"}

const (
	ErrFailedOpenConfigFile  = "failed to open configFil"
	ErrFailedReadConfig      = "failed to read config"
	ErrMissingDSN            = "missing DSN"
	ErrWrongSQLDriver        = "wrong SQL driver"
	ErrMissingOutputPaths    = "missing logger output paths"
	ErrMissingErrOutputPaths = "missing logger error output paths"
)
