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

type GetMangaChapterService struct {
}

type GetMangaChapterRequest struct {
	Id      string `json:"id"`
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`

	// TODO: Decouple from amq dependencies and publisher.
	Channel   *amqp091.Channel
	Publisher events.Publisher
}
type GetMangaChapterResponse struct {
	Chapter entities.Chapter
}

func (self *GetMangaChapterService) Exec(request GetMangaChapterRequest) (*GetMangaChapterResponse, error) {
	reply_queue, err := channels.CreateReplyQueue(request.Channel)
	if err != nil {
		log.Printf("[CHAPTER] - Unable to create reply queue, error %+v\n", err)
		return nil, err
	}

	message := events.ScrapeChapterMessage{
		Slug:    request.Slug,
		Chapter: request.Chapter,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id:    request.Id,
				Reply: reply_queue.Name,
			},
		},
	}

	_, err = request.Publisher(message)
	if err != nil {
		log.Printf("[CHAPTER] - Unable to publish message, error %+v\n", err)
		return nil, err
	}

	reply_message, err := request.Channel.Consume(reply_queue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Printf("[PUBLISHER] - Unable to consume reply queue messages, error %+v\n", err)
		return nil, err
	}

	log.Println("[PUBLISHER] - Waiting for reply...")
	select {
	case reply := <-reply_message:
		{
			log.Println("[PUBLISHER] - Message received")
			message := events.ScrapedChapterMessage{}

			err := json.Unmarshal(reply.Body, &message)
			if err != nil {
				log.Printf("[PUBLISHER] - Unable to parse reply body, error %+v\n", err)
				return nil, err
			}

			response := GetMangaChapterResponse{
				Chapter: message.Data,
			}

			return &response, nil
		}
	case <-time.After(time.Duration(channels.QUEUE_TTL) * time.Second):
		{
			return nil, errors.New("Unable to retrieve data, request timeout")
		}
	}

}
