package config

type Config struct {
	Logger  ConfigLogger  `json:"logger"`
	Server  ConfigServer  `json:"server"`
	Storage ConfigStorage `json:"storage"`
}

type ConfigLogger struct {
	Level            string   `json:"level"`
	OutputPaths      []string `json:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths"`
}

type ConfigServer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ConfigStorage struct {
	Type      string `json:"type"`
	SQLDriver string `json:"sql_driver"`
	DSN       string `json:"dsn"`
}
