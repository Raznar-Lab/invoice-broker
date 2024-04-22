package config

type MessageBrokerConfig struct {
	RabbitMQ RabbitMQ `yaml:"rabbit_mq" validate:"required"`
}

type RabbitMQ struct {
	ExchangeKey string `yaml:"exchange_key" validate:"required"`
	Host        string `yaml:"host" validate:"required"`
	Port        string `yaml:"port" validate:"required"`
	Username    string `yaml:"username" validate:"required"`
	Password    string `yaml:"password" validate:"required"`
}
