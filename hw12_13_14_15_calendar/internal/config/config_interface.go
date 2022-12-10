package config

type Config struct {
	Logger  Logger  `json:"logger"`
	Server  Server  `json:"server"`
	Storage Storage `json:"storage"`
}

type Logger struct {
	Level            string   `json:"level"`
	OutputPaths      []string `json:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Storage struct {
	Type      string `json:"type"`
	SQLDriver string `json:"sql_driver"`
	DSN       string `json:"dsn"`
}
