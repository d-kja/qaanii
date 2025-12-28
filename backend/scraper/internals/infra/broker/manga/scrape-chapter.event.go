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

func ScrapeChapterSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	scraped_chapter_publisher, ok := ctx.Value(events.SCRAPED_CHAPTER_EVENT).(*events.Publisher)
	if !ok {
		log.Println("Scrape chapter consumer | Unable to retrieve scraped chapter publisher.")
		return errors.New("Unable to retrieve scrape chapter publisher")
	}

	message := events.ScrapeChapterMessage{}
	err := json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("Scrape chapter consumer | unable to parse message body, error: %+v\n", err)
		return err
	}

	slug_len := len(message.Slug)
	if slug_len == 0 {
		log.Println("Scrape chapter consumer | missing slug parameter")
		return errors.New("Missing slug parameter")
	}

	chapter_len := len(message.Chapter)
	if chapter_len == 0 {
		log.Println("Scrape chapter consumer | missing chapter parameter")
		return errors.New("Missing chapter parameter")
	}

	service := usecase.GetMangaChapterService{
		Scraper: coreentities.NewScraper(),
	}

	response, err := service.Exec(usecase.GetMangaChapterRequest{
		Slug:    message.Slug,
		Chapter: message.Chapter,
	})

	pub_message := events.ScrapedChapterMessage{
		BaseEvent: events.BaseEvent{},
		Data:      response.Chapter,
	}

	pub_message.GenerateEventId(string(events.SCRAPED_CHAPTER_EVENT), "n/a")

	_, err = (*scraped_chapter_publisher)(pub_message)
	if err != nil {
		log.Printf("Scrape manga publisher | Unable to publish scraped results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}

	return err
}
