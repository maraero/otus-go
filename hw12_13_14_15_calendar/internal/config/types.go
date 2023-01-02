package config

import "time"

type CalendarConfig struct {
	Logger  Logger  `json:"logger"`
	Server  Server  `json:"server"`
	Storage Storage `json:"storage"`
}

type SchedulerConfig struct {
	Logger Logger `json:"logger"`
	RMQ    RMQ    `json:"rmq"`
}

type SenderConfig struct {
	Logger Logger `json:"logger"`
	RMQ    RMQ    `json:"rmq"`
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

type RMQ struct {
	URI          string    `json:"uri"`
	Exchange     string    `json:"exchange"`
	ExchangeType string    `json:"exchangeType"`
	Queue        string    `json:"queue"`
	BindingKey   string    `json:"bindingKey"`
	ConsumerTag  string    `json:"consumerTag"`
	Lifetime     time.Time `json:"lifetime"`
}
