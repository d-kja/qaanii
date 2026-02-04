package broker

import (
	"context"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SetupSubscribers(request broker.SubscriberRequest) {
	sample := func(queue amqp.Delivery, ctx *context.Context) error {
		// Temporary mock function to replace copied logic
		return nil
	}

	broker.CreateConsumer(events.SCRAPED_CHAPTER_EVENT, request, sample)
	broker.CreateConsumer(events.SCRAPED_MANGA_EVENT, request, sample)
	broker.CreateConsumer(events.SEARCHED_MANGA_EVENT, request, sample)
}

