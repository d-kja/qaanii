package manga

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/mangas/use_case"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ScrapeMangaSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	scraped_manga_publisher, ok := ctx.Value(events.SCRAPED_MANGA_EVENT).(*events.Publisher)
	if !ok {
		log.Println("[SUBSCRIBER/MANGA] - Unable to retrieve scraped manga publisher.")
		return errors.New("Unable to retrieve scrape manga publisher")
	}

	message := events.ScrapeMangaMessage{}
	err := json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - unable to parse message body, error: %+v\n", err)
		return err
	}

	slug_len := len(message.Slug)
	if slug_len == 0 {
		log.Println("[SUBSCRIBER/MANGA] - missing slug parameter")
		return errors.New("Missing slug parameter")
	}

	service := usecase.GetMangaBySlugService{
		Scraper: coreentities.NewScraper(),
	}

	response, err := service.Exec(usecase.GetMangaBySlugRequest{
		Slug: message.Slug,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - Unable to retrieve scraped results, error: %+v\n", err)
		return errors.New("Unable to retrieve results")
	}

	pub_message := events.ScrapedMangaMessage{
		BaseEvent: events.BaseEvent{},
		Data:      response.Manga,
	}

	pub_message.GenerateEventId(string(events.SCRAPED_MANGA_EVENT), "n/a")

	_, err = (*scraped_manga_publisher)(pub_message)
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - Unable to publish scraped results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}

	return err
}
