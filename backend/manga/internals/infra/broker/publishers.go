package broker

import (
	"context"
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

func SetupPublishers(request PublisherRequest) {
	create_publisher(events.SCRAPE_CHAPTER_EVENT, request)
	create_publisher(events.SCRAPE_MANGA_EVENT, request)

	create_publisher(events.SEARCH_MANGA_EVENT, request)
}

func create_publisher(event events.Events, request PublisherRequest) {
	queue, err := channels.CreateQueue(string(event), request.Channel)
	if err != nil {
		log.Panicf("Broker publisher | unable to create queue [%v], error: %+v\n", event, err)
		return
	}

	var handler events.Publisher = func(data any) (any, error) {
		response, err := channels.PublishMessage(data, queue, request.Channel)
		if err != nil {
			log.Printf("Broker publisher | unable to publish message, error: %+v\n", err)
			return nil, err
		}

		log.Printf("Broker publisher | %v message sent\n", event)
		return response, nil
	}

	log.Printf("[%v] Publisher created \n", event)
	*request.Context = context.WithValue(*request.Context, event, &handler)
}
