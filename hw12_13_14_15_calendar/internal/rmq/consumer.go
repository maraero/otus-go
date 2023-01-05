package rmq

import (
	"context"
	"fmt"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type RMQConsumer struct {
	rmqc        *RMQConnection
	channel     *rabbitmq.Channel
	consumerTag string
	queueName   string
}

type Worker func(ctx context.Context, d amqp.Delivery) error

func NewRMQConsumer(
	ctx context.Context,
	rmqc *RMQConnection,
	exchangeName,
	exchangeType,
	queueName,
	bindingKey,
	consumerTag string,
	logger *logger.Log,
	worker Worker,
) (*RMQConsumer, error) {
	channel, err := rmqc.OpenChannel()
	if err != nil {
		return nil, err
	}

	if err = declareExchange(channel, exchangeName, exchangeType); err != nil {
		return nil, err
	}

	if err = declareQueue(channel, queueName); err != nil {
		return nil, err
	}

	if err = bindQueue(channel, exchangeName, queueName, bindingKey); err != nil {
		return nil, err
	}

	deliveries, err := channel.Consume(
		queueName,   // name
		consumerTag, // consumerTag,
		false,       // noAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("can not consume message: %w", err)
	}

	go handle(ctx, deliveries, logger, worker)

	return &RMQConsumer{rmqc: rmqc, consumerTag: consumerTag, queueName: queueName}, nil
}

func handle(ctx context.Context, deliveries <-chan amqp.Delivery, logger *logger.Log, worker Worker) {
	for d := range deliveries {
		err := worker(ctx, d)
		if err != nil {
			logger.Error("can not handle RMQ message: %w", err)
			err = d.Ack(false)
			if err != nil {
				logger.Error("can not nAck RMQ message: %w", err)
			}
		} else {
			err = d.Ack(true)
			if err != nil {
				logger.Error("can not Ack RMQ message: %w", err)
			}
		}
	}
}

func declareQueue(channel *rabbitmq.Channel, queueName string) error {
	if _, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	); err != nil {
		return fmt.Errorf("can not declare queue: %w", err)
	}
	return nil
}

func bindQueue(channel *rabbitmq.Channel, exchangeName, queueName, bindingKey string) error {
	if err := channel.QueueBind(queueName, bindingKey, exchangeName, false, nil); err != nil {
		return fmt.Errorf("can not bind queue: %w", err)
	}
	return nil
}
