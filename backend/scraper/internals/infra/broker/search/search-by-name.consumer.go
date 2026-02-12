package search

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	coreutils "qaanii/scraper/internals/domain/core/core_utils"
	usecase "qaanii/scraper/internals/domain/search/use_case"
	internal_utils "qaanii/scraper/internals/utils"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SearchByNameSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	pool, ok := ctx.Value(internal_utils.SCRAPER_POOL_KEY).(*coreutils.BrowserPool)
	if !ok {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to retrieve scraper pool.")
		return errors.New("Unable to retrieve scraper pool")
	}

	connection, ok := ctx.Value(internal_utils.QUEUE_CONNECTION_KEY).(*amqp.Connection)
	if !ok {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to retrieve base connection.")
		return errors.New("Unable to retrieve base connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - Unable to retrieve base queue, error: %v", err)
		return errors.Join(errors.New("Unable to create queue"), err)
	}
	defer channel.Close()

	message := events.SearchMangaMessage{}
	err = json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - unable to parse message body, error: %+v\n", err)
		return err
	}

	query_len := len(message.Query)
	if query_len == 0 {
		log.Println("[SUBSCRIBER/SEARCH] - missing query parameter")
		return errors.New("Missing query parameter")
	}

	scraper, err := pool.Get()
	if err != nil || scraper == nil {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to acquire scraper instance")
		return errors.New("Unable to acquire scraper instance")
	}
	defer pool.Release(scraper)

	service := usecase.SearchByNameService{
		Scraper: *scraper,
	}

	response, err := service.Exec(usecase.SearchByNameRequest{
		Search: message.Query,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - Unable to retrieve search results, error: %+v\n", err)
		return errors.New("Unable to retrieve search results")
	}

	pub_message := events.SearchedMangaMessage{
		BaseEvent: events.BaseEvent{
			Metadata: message.Metadata,
		},
		Data: response.Mangas,
	}

	_, err = broker.Reply(message.Metadata.Reply, broker.PublishRequest{
		Channel:    channel,
		Connection: connection,
		Data:       pub_message,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - Unable to publish search results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}

	return err
}
