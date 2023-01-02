package config

type CalendarConfig struct {
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
	Host     string `json:"host"`
	HTTPPort string `json:"httpPort"`
	GrpcPort string `json:"grpcPort"`
}

type Storage struct {
	Type   string `json:"type"`
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}
