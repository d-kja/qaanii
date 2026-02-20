package manga

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	coreutils "qaanii/scraper/internals/domain/core/core_utils"
	usecase "qaanii/scraper/internals/domain/mangas/use_case"
	internal_utils "qaanii/scraper/internals/utils"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ScrapeChapterSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	pool, ok := ctx.Value(internal_utils.SCRAPER_POOL_KEY).(*coreutils.BrowserPool)
	if !ok {
		log.Println("[SUBSCRIBER/CHAPTER] - Unable to retrieve scraper pool.")
		return errors.New("Unable to retrieve scraper pool")
	}

	connection, ok := ctx.Value(internal_utils.QUEUE_CONNECTION_KEY).(*amqp.Connection)
	if !ok {
		log.Println("[SUBSCRIBER/CHAPTER] - Unable to retrieve base connection.")
		return errors.New("Unable to retrieve base connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("[SUBSCRIBER/CHAPTER]- Unable to retrieve base queue, error: %v", err)
		return errors.Join(errors.New("Unable to create queue"), err)
	}
	defer channel.Close()

	message := events.ScrapeChapterMessage{}
	err = json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("[SUBSCRIBER/CHAPTER] - unable to parse message body, error: %+v\n", err)
		return err
	}

	slug_len := len(message.Slug)
	if slug_len == 0 {
		log.Println("[SUBSCRIBER/CHAPTER] - missing slug parameter")
		return errors.New("Missing slug parameter")
	}

	chapter_len := len(message.Chapter)
	if chapter_len == 0 {
		log.Println("[SUBSCRIBER/CHAPTER] - missing chapter parameter")
		return errors.New("Missing chapter parameter")
	}

	scraper, err := pool.Get()
	if err != nil || scraper == nil {
		log.Println("[SUBSCRIBER/CHAPTER] - Unable to acquire scraper instance")
		return errors.New("Unable to acquire scraper instance")
	}
	defer pool.Release(scraper)

	service := usecase.GetMangaChapterService{
		Scraper: *scraper,
	}

	response, err := service.Exec(usecase.GetMangaChapterRequest{
		Slug:    message.Slug,
		Chapter: message.Chapter,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/CHAPTER] - Unable to retrieve chapter result, error: %+v\n", err)
		return errors.New("Unable to retrieve chapter")
	}

	pages := *response.Chapter.Pages

	pub_message := events.ScrapedChapterMessage{
		BaseEvent: events.BaseEvent{
			Metadata: message.Metadata,
		},
		Data: events.MessageChapter{
			Title: response.Chapter.Title,
			Link:  response.Chapter.Link,
			Time:  response.Chapter.Time,
			Pages: pages,
		},
	}

	a, _ := json.MarshalIndent(pub_message, "", " ")
	log.Printf("JSON: %+v", string(a))

	_, err = broker.Reply(message.Metadata.Reply, broker.PublishRequest{
		Channel:    channel,
		Connection: connection,
		Data:       pub_message,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/CHAPTER] - Unable to publish chapter results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}

	return err
}
