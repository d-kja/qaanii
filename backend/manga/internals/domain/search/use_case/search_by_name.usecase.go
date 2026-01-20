package usecase

import (
	// "context"
	// "fmt"
	// "net/url"
	"qaanii/shared/entities"
	// "qaanii/shared/utils"

	// amq "github.com/rabbitmq/amqp091-go"
)

type SearchByNameService struct {
}

type SearchByNameRequest struct {
	Search string
	// Context *context.Context
}
type SearchByNameResponse struct {
	Mangas []entities.Manga
}

func (self *SearchByNameService) Exec(request SearchByNameRequest) (*SearchByNameResponse, error) {
	// envs := utils.Utils{}.Envs()

	// ctx := *request.Context
	// search_manga_publisher, ok := ctx.Value()

	// search := url.QueryEscape(request.Search)
	// url := fmt.Sprintf("%v/search?q=%v", envs["base_url"], search)

	mangas := []entities.Manga{}

	// Send event to message broker
	
	// amq.

	// Lock thread waiting for channel response with timeout logic.

	response := SearchByNameResponse{
		Mangas: mangas, 
	}

	return &response, nil
}
