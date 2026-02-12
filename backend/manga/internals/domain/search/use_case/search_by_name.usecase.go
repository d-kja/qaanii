package usecase

import (
	"encoding/json"
	"errors"
	"log"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"qaanii/shared/entities"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type SearchByNameService struct {
}

type SearchByNameRequest struct {
	Id        string
	Search    string
	Channel   *amqp091.Channel
	Publisher events.Publisher
}
type SearchByNameResponse struct {
	Mangas []entities.Manga
}

var err_response base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

func (self *SearchByNameService) Exec(request SearchByNameRequest) (*SearchByNameResponse, error) {
	reply_queue, err := channels.CreateReplyQueue(request.Channel)
	if err != nil {
		log.Printf("[SEARCH] - Unable to create reply queue, error %+v\n", err)
		return nil, err
	}

	message := events.SearchMangaMessage{
		Query: request.Search,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id:    request.Id,
				Reply: reply_queue.Name,
			},
		},
	}

	request.Publisher(message)
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
			message := events.SearchedMangaMessage{}

			err := json.Unmarshal(reply.Body, &message)
			if err != nil {
				log.Printf("[PUBLISHER] - Unable to parse reply body, error %+v\n", err)
				return nil, err
			}

			response := SearchByNameResponse{
				Mangas: message.Data,
			}

			return &response, nil
		}

	case <-time.After(120 * time.Second):
		{
			return nil, errors.New("Unable to retrieve data, request timeout")
		}
	}
}
