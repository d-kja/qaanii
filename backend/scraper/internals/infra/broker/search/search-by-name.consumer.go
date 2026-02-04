package search

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/search/use_case"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SearchByNameSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	searched_publisher, ok := ctx.Value(events.SEARCHED_MANGA_EVENT).(*events.Publisher)
	if !ok {
		log.Println("[SUBSCRIBER/SEARCH] - Unable to retrieve searched manga publisher.")
		return errors.New("Unable to retrieve search results publisher")
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

	pub_message := events.SearchedMangaMessage{
		BaseEvent: events.BaseEvent{},
		Data: response.Mangas,
	}

	pub_message.GenerateEventId(string(events.SEARCHED_MANGA_EVENT), "n/a")

	_, err = (*searched_publisher)(pub_message)
	if err != nil {
		log.Printf("[SUBSCRIBER/SEARCH] - Unable to publish search results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}

	return err
}
