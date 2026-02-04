package broker

import (
	"context"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SetupSubscribers(request broker.SubscriberRequest) {
	// Not necessary, but to avoid sync issues during hot-reload I added a wg
	wg := sync.WaitGroup{}
	wg.Add(3)

	sample := func(queue amqp.Delivery, ctx *context.Context) error {
		// Temporary mock function to replace copied logic
		return nil
	}

	go broker.CreateConsumer(events.SCRAPED_CHAPTER_EVENT, request, sample, &wg)
	go broker.CreateConsumer(events.SCRAPED_MANGA_EVENT, request, sample, &wg)
	go broker.CreateConsumer(events.SEARCHED_MANGA_EVENT, request, sample, &wg)

	wg.Wait()
}

