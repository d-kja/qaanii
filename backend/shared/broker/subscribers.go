package broker

import (
	"context"
	"encoding/json"
	"log"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberRequest struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Context    *context.Context
}

func CreateConsumer(queue events.Events, request SubscriberRequest, callback func(amqp.Delivery, *context.Context) error) {
	queue_ch, err := channels.CreateQueue(string(queue), request.Channel)

	if err != nil {
		log.Printf("[BROKER/SUBSCRIBER] - Manga queue creation failed, error: %+v\n", err)
		return
	}

	messages_ch, err := channels.CreateSubscriberChannel(queue_ch, request.Channel)
	if err != nil {
		log.Printf("[BROKER/SUBSCRIBER] - Unable to setup consumer for %v, error: %+v\n", queue, err)
		return
	}

	go HandleMessages(messages_ch, request.Context, callback)
	log.Printf("[%v] - subscriber created\n", queue)
}

func HandleMessages(messages_ch <-chan amqp.Delivery, ctx *context.Context, callback func(amqp.Delivery, *context.Context) error) {
	for raw_message := range messages_ch {
		message := events.BaseEvent{}

		// Validate initial event
		err := json.Unmarshal(raw_message.Body, &message)
		if err != nil {
			log.Printf("[BROKER/SUBSCRIBER] - Unable to parse message body, \n - error: %+v - payload: %v\n", err, string(raw_message.Body))
			continue
		}

		if len(message.Metadata.Id) == 0 {
			log.Printf("[BROKER/SUBSCRIBER] - Invalid event, skipping: %v\n", string(raw_message.Body))

			err = raw_message.Nack(false, false) // Invalid event, we can't process it.
			if err != nil {
				log.Printf("[BROKER/SUBSCRIBER] - Unable to acknowledge message, error: %+v\n", err)
			}

			continue
		}

		err = callback(raw_message, ctx)
		if err != nil {
			log.Printf("[BROKER/SUBSCRIBER] - unable to process message, error: %+v\n", err)

			err = raw_message.Nack(false, true)
			if err != nil {
				log.Printf("[BROKER/SUBSCRIBER] - Unable to acknowledge message, error: %+v\n", err)
			}

			continue
		}

		err = raw_message.Ack(false)
		if err != nil {
			log.Printf("[BROKER/SUBSCRIBER] - Unable to acknowledge message, error: %+v\n", err)
		}
	}
}
