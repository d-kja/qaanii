package search

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/search/use_case"
	internal_utils "qaanii/scraper/internals/utils"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SearchByNameSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	connection, ok := ctx.Value(internal_utils.QUEUE_CONNECTION_KEY).(*amqp.Connection)
	if !ok {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to retrieve base connection.")
		return errors.New("Unable to retrieve base connection")
	}

	channel, ok := ctx.Value(internal_utils.QUEUE_CHANNEL_KEY).(*amqp.Channel)
	if !ok {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to retrieve base queue.")
		return errors.New("Unable to retrieve base queue")
	}

	message := events.SearchMangaMessage{}
	err := json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - unable to parse message body, error: %+v\n", err)
		return err
	}

	query_len := len(message.Query)
	if query_len == 0 {
		log.Println("[SUBSCRIBER/SEARCH] - missing query parameter")
		return errors.New("Missing query parameter")
	}

	service := usecase.SearchByNameService{
		Scraper: coreentities.NewScraper(),
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
