package usecase

import (
	"fmt"
	"net/url"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type SearchByNameService struct {
}

type SearchByNameRequest struct {
	Search string
}
type SearchByNameResponse struct {
	Mangas []entities.Manga
}

func (self *SearchByNameService) Exec(request SearchByNameRequest) (*SearchByNameResponse, error) {
	envs := utils.Utils{}.Envs()

	search := url.QueryEscape(request.Search)
	_ = fmt.Sprintf("%v/search?q=%v", envs["base_url"], search)

	response := SearchByNameResponse{}

	return &response, nil
}
