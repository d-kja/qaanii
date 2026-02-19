package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"qaanii/shared/entities"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type GetMangaBySlugService struct {
}

type GetMangaBySlugRequest struct {
	Id   string
	Slug string `json:"slug"`

	Channel   *amqp091.Channel
	Publisher events.Publisher
}
type GetMangaBySlugResponse struct {
	Manga entities.Manga
}

func (self *GetMangaBySlugService) Exec(request GetMangaBySlugRequest) (*GetMangaBySlugResponse, error) {
	reply_queue, err := channels.CreateReplyQueue(request.Channel)
	if err != nil {
		log.Printf("[MANGA] - Unable to create reply queue, error %+v\n", err)
		return nil, err
	}

	message := events.ScrapeMangaMessage{
		Slug: request.Slug,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id:    request.Id,
				Reply: reply_queue.Name,
			},
		},
	}

	_, err = request.Publisher(message)
	if err != nil {
		log.Printf("[MANGA] - Unable to publish message, error %+v\n", err)
		return nil, err
	}

	reply_message, err := request.Channel.Consume(reply_queue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Printf("[MANGA] - Unable to create consumer, error %+v\n", err)
		return nil, err
	}

	log.Println("[PUBLISHER] - Waiting for reply...")
	select {
	case reply := <-reply_message:
		{
			log.Println("[PUBLISHER] - Message received")
			message := events.ScrapedMangaMessage{}

			err := json.Unmarshal(reply.Body, &message)
			if err != nil {
				log.Printf("[MANGA] - Unable to create consumer, error %+v\n", err)
				return nil, err
			}

			response := GetMangaBySlugResponse{
				Manga: message.Data,
			}

			return &response, nil
		}
	case <-time.After(time.Duration(channels.QUEUE_TTL) * time.Second):
		{
			return nil, errors.New("Unable to retrieve data, request timeout")
		}
	}
}
