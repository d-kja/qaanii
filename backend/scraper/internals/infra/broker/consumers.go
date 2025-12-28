package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerRequest struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func SetupConsumers(request ConsumerRequest) {
	debug_fn := func(raw_message amqp.Delivery) error {
		tmp := map[string]any{}
		err := json.Unmarshal(raw_message.Body, &tmp)
		if err != nil {
			log.Printf("X Consumer | unable to parse message body, error: %+v\n", err)
			return err
		}

		marshaled, err := json.MarshalIndent(tmp, "", "   ")
		if err != nil {
			log.Printf("marshaling error: %s\n", err)
		}

		fmt.Println(string(marshaled))
		return nil
	}

	// Search consumer
	create_consumer(events.SEARCH_MANGA_EVENT, debug_fn, request)
}

func create_consumer(queue events.Events, callback func(amqp.Delivery) error, request ConsumerRequest) error {
	queue_ch, err := channels.CreateQueue(string(queue), request.Channel)
	if err != nil {
		log.Panicf("Consumer | Manga queue creation failed, error: %+v", err)

		return err
	}

	messages_ch, err := channels.CreateConsumerChannel(queue_ch, request.Channel)
	if err != nil {
		log.Panicf("Consumer | Unable to setup consumer for %v, error: %+v\n", queue, err)
	}

	lock_ch := make(chan bool)

	go handle_messages(messages_ch, callback)
	log.Println("Consumer | Waiting for queue messages")

	<-lock_ch
	return nil
}

func handle_messages(messages_ch <-chan amqp.Delivery, callback func(amqp.Delivery) error) {
	for raw_message := range messages_ch {
		message := events.BaseEvent{}

		// Validate initial event
		err := json.Unmarshal(raw_message.Body, &message)
		if err != nil {
			log.Printf("Consumer | unable to parse message body, error: %+v\n", err)
			continue
		}

		err = callback(raw_message)
		if err != nil {
			log.Printf("Consumer | unable to process message, error: %+v\n", err)
			continue
		}

		err = raw_message.Ack(false)
		if err != nil {
			log.Printf("Consumer | Unable to acknowledge message, error: %+v\n", err)
		}
	}
}
