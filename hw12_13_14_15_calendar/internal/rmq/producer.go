package rmq

import (
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
)

type Producer struct {
	rmqc         *Connection
	channel      *rabbitmq.Channel
	exchangeName string
}

func NewRMQProducer(rmqc *Connection, exchangeName, exchangeType string) (*Producer, error) {
	channel, err := rmqc.OpenChannel()
	if err != nil {
		return nil, err
	}

	if err = declareExchange(channel, exchangeName, exchangeType); err != nil {
		return nil, err
	}

	return &Producer{rmqc: rmqc, channel: channel, exchangeName: exchangeName}, nil
}

func (rmqp *Producer) Publish(key string, body string) error {
	return rmqp.channel.Publish(rmqp.exchangeName, key, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		ContentEncoding: "application/json",
		Body:            []byte(body),
	})
}
