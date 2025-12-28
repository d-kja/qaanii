package channels

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MANGA_CHANNEL string = "@manga"
)

func CreateQueue(name string, channel *amqp.Channel) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		name,
		true,  // Durable
		false, // Delete when used
		false, // Exclusive
		false, // No-wait
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
