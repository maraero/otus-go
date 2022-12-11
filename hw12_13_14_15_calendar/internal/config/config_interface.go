package config

type Config struct {
	Logger  Logger  `json:"logger"`
	Server  Server  `json:"server"`
	Storage Storage `json:"storage"`
}

type Logger struct {
	Level            string   `json:"level"`
	OutputPaths      []string `json:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Storage struct {
	Type     string `json:"type"`
	Database string `json:"database"`
	DSN      string `json:"dsn"`
}
