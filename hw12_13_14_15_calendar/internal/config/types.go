package config

type CalendarConfig struct {
	Logger  Logger  `json:"logger"`
	Server  Server  `json:"server"`
	Storage Storage `json:"storage"`
}

type SchedulerConfig struct {
	Logger Logger       `json:"logger"`
	RMQ    RMQScheduler `json:"rmq"`
}

type SenderConfig struct {
	Logger Logger    `json:"logger"`
	RMQ    RMQSender `json:"rmq"`
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
	URI          string `json:"uri"`
	Exchange     string `json:"exchange"`
	ExchangeType string `json:"exchangeType"`
}

type RMQScheduler struct {
	RMQ
}

type RMQSender struct {
	RMQ
	Queue       string `json:"queue"`
	ConsumerTag string `json:"consumerTag"`
}
