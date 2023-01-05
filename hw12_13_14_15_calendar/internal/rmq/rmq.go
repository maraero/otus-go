package rmq

import (
	"fmt"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type RMQConnection struct {
	conn        *rabbitmq.Connection
	channelList []*rabbitmq.Channel
	logger      *logger.Log
}

func NewRMQConnection(uri string, lggr *logger.Log) (*RMQConnection, error) {
	conn, err := rabbitmq.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("can not set up RMQ connection: %w", err)
	}
	return &RMQConnection{conn: conn}, nil
}

func (rmqc *RMQConnection) OpenChannel() (*rabbitmq.Channel, error) {
	channel, err := rmqc.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can not open RMQ channel: %w", err)
	}
	rmqc.channelList = append(rmqc.channelList, channel)
	return channel, nil
}

func (rmqc *RMQConnection) Shutdown() {
	for _, ch := range rmqc.channelList {
		if err := ch.Close(); err != nil {
			rmqc.logger.Error("can not close RMQ channel")
		}
	}

	if err := rmqc.conn.Close(); err != nil {
		rmqc.logger.Error("can not close RMQ connection")
	}
}

func declareExchange(channel *rabbitmq.Channel, exchangeName, exchangeType string) error {
	if err := channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("can not declare RMQ exchange: %w", err)
	}

	return nil
}
