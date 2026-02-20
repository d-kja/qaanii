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

func ScrapeMangaSubscriber(raw_message amqp.Delivery, ctx_prt *context.Context) error {
	ctx := *ctx_prt

	pool, ok := ctx.Value(internal_utils.SCRAPER_POOL_KEY).(*coreutils.BrowserPool)
	if !ok {
		log.Println("[SUBSCRIBER/MANGA] - Unable to retrieve scraper pool.")
		return errors.New("Unable to retrieve scraper pool")
	}

	connection, ok := ctx.Value(internal_utils.QUEUE_CONNECTION_KEY).(*amqp.Connection)
	if !ok {
		log.Println("[SUBSCRIBER/MANGA] - Unable to retrieve base connection.")
		return errors.New("Unable to retrieve base connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - Unable to retrieve base queue, error: %v", err)
		return errors.Join(errors.New("Unable to create queue"), err)
	}
	defer channel.Close()

	message := events.ScrapeMangaMessage{}
	err = json.Unmarshal(raw_message.Body, &message)
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - unable to parse message body, error: %+v\n", err)
		return err
	}

	slug_len := len(message.Slug)
	if slug_len == 0 {
		log.Println("[SUBSCRIBER/MANGA] - missing slug parameter")
		return errors.New("Missing slug parameter")
	}

	scraper, err := pool.Get()
	if err != nil || scraper == nil {
		log.Println("[SUBSCRIBER/MANGA] - Unable to acquire scraper instance")
		return errors.New("Unable to acquire scraper instance")
	}
	defer pool.Release(scraper)

	service := usecase.GetMangaBySlugService{
		Scraper: *scraper,
	}

	response, err := service.Exec(usecase.GetMangaBySlugRequest{
		Slug: message.Slug,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - Unable to retrieve scraped results, error: %+v\n", err)
		return errors.New("Unable to retrieve results")
	}

	pub_message := events.ScrapedMangaMessage{
		BaseEvent: events.BaseEvent{
			Metadata: message.Metadata,
		},
		Data:      response.Manga,
	}

	_, err = broker.Reply(message.Metadata.Reply, broker.PublishRequest{
		Channel:    channel,
		Connection: connection,
		Data:       pub_message,
	})
	if err != nil {
		log.Printf("[SUBSCRIBER/MANGA] - Unable to publish manga results, error: %+v\n", err)
		return errors.New("Unable to publish results")
	}


	return err
}
