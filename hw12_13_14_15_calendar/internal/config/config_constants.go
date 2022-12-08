package config

const (
	StorageInMemory string = "in-memory"
	StorageSql      string = "SQL"
)

const (
	ErrFailedOpenConfigFile  = "failed to open configFil"
	ErrFailedReadConfig      = "failed to read config"
	ErrMissingDSN            = "missing DSN"
	ErrMissingOutputPaths    = "missing logger output paths"
	ErrMissingErrOutputPaths = "missing logger error output paths"
)
