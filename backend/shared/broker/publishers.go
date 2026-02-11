package broker

import (
	"context"
	"encoding/json"
	"log"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherRequest struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Context    *context.Context
}

func CreatePublisher(event events.Events, request PublisherRequest) {
	queue, err := channels.CreateQueue(string(event), request.Channel)
	if err != nil {
		log.Printf("[BROKER/PUBLISHER] - unable to create queue [%v], error: %+v\n", event, err)
		return
	}

	var handler events.Publisher = func(data any) (any, error) {
		response, err := channels.PublishMessage(data, queue, request.Channel)
		if err != nil {
			log.Printf("[BROKER/PUBLISHER] - unable to publish message, error: %+v\n", err)
			return nil, err
		}

		log.Printf("[BROKER/PUBLISHER] - %v message sent\n", event)
		return response, nil
	}

	log.Printf("[%v] - publisher created, queue\n", event)
	*request.Context = context.WithValue(*request.Context, event, &handler)
}

type PublishRequest struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Data       any
}

func Reply(queue string, request PublishRequest) (any, error) {
	body, err := json.Marshal(request.Data)
	if err != nil {
		log.Printf("[BROKER/PUBLISHER] - unable to marshal message, error: %+v\n", err)
		return nil, err
	}

	err = request.Channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("[BROKER/PUBLISHER] - unable to publish message, error: %+v\n", err)
		return nil, err
	}

	log.Printf("[BROKER/PUBLISHER] - %v message sent\n", queue)
	return request.Data, nil
}
