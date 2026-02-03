package broker

import (
	"context"
	"qaanii/manga/internals/utils"
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
	l := utils.GetLogger()

	queue, err := channels.CreateQueue(string(event), request.Channel)
	if err != nil {
		l.Panicf("[BROKER/PUBLISHER] - unable to create queue [%v], error: %+v", event, err)
		return
	}

	var handler events.Publisher = func(data any) (any, error) {
		response, err := channels.PublishMessage(data, queue, request.Channel)
		if err != nil {
			l.Infof("[BROKER/PUBLISHER] - unable to publish message, error: %+v", err)
			return nil, err
		}

		l.Infof("[BROKER/PUBLISHER] - %v message sent", event)
		return response, nil
	}

	l.Infof("[%v] - publisher created, queue", event)
	*request.Context = context.WithValue(*request.Context, event, &handler)
}
