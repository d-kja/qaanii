package channels

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MANGA_CHANNEL string = "@manga"
	QUEUE_TTL     uint    = 120
)

func CreateQueue(name string, channel *amqp.Channel) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		name,
		false, // Durable
		true,  // Delete when used
		false, // Exclusive
		false, // No-wait
		amqp.Table{
			"x-message-ttl": QUEUE_TTL * 100, // messages expire after 120s
			"x-max-length":  1000,            // max 1000 messages, oldest dropped
			"x-overflow":    "drop-head",     // drop oldest when full
		},
	)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}

// Single run queue, for async work mixed with grpc stream...
func CreateReplyQueue(channel *amqp.Channel) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		"", // Generate random name
		false,
		false,
		true, // Exclusive
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func PublishMessage(data any, queue *amqp.Queue, channel *amqp.Channel) (any, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateSubscriberChannel(queue *amqp.Queue, channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	messages_ch, err := channel.Consume(
		queue.Name,
		"",    // Consumer
		false, // Auto acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messages_ch, nil
}
