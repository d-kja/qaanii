package broker

import (
	"context"
	"encoding/json"
	"log"
	"qaanii/scraper/internals/infra/broker/manga"
	"qaanii/scraper/internals/infra/broker/search"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberRequest struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Context    *context.Context
}

func SetupSubscribers(request SubscriberRequest) {
	// Not necessary, but to avoid sync issues during hot-reload I added a wg
	wg := sync.WaitGroup{}
	wg.Add(3)

	go create_consumer(events.SCRAPED_CHAPTER_EVENT, request, manga.ScrapeMangaSubscriber, &wg)
	go create_consumer(events.SCRAPED_MANGA_EVENT, request, manga.ScrapeChapterSubscriber, &wg)

	go create_consumer(events.SEARCHED_MANGA_EVENT, request, search.SearchByNameSubscriber, &wg)

	wg.Wait()
}

func create_consumer(queue events.Events, request SubscriberRequest, callback func(amqp.Delivery, *context.Context) error, wg *sync.WaitGroup) {
	queue_ch, err := channels.CreateQueue(string(queue), request.Channel)
	if err != nil {
		log.Panicf("Subscriber | Manga queue creation failed, error: %+v", err)
		return
	}

	messages_ch, err := channels.CreateSubscriberChannel(queue_ch, request.Channel)
	if err != nil {
		log.Panicf("Subscriber | Unable to setup consumer for %v, error: %+v\n", queue, err)
		return
	}

	lock_ch := make(chan bool)

	go handle_messages(messages_ch, request.Context, callback)
	log.Printf("[%v] Subscriber created \n", queue)

	wg.Done()

	<-lock_ch
}

func handle_messages(messages_ch <-chan amqp.Delivery, ctx *context.Context, callback func(amqp.Delivery, *context.Context) error) {
	for raw_message := range messages_ch {
		message := events.BaseEvent{}

		// Validate initial event
		err := json.Unmarshal(raw_message.Body, &message)
		if err != nil {
			log.Printf("Subscriber | Unable to parse message body, \n - error: %+v\n - payload: %v\n", err, string(raw_message.Body))
			continue
		}

		if len(message.Metadata.Id) == 0 {
			log.Printf("Subscriber | Invalid event, skipping: %v\n", string(raw_message.Body))

			err = raw_message.Ack(false)
			if err != nil {
				log.Printf("Subscriber | Unable to acknowledge message, error: %+v\n", err)
			}

			continue
		}

		err = callback(raw_message, ctx)
		if err != nil {
			log.Printf("Subscriber | unable to process message, error: %+v\n", err)
			continue
		}

		err = raw_message.Ack(false)
		if err != nil {
			log.Printf("Subscriber | Unable to acknowledge message, error: %+v\n", err)
		}
	}
}
