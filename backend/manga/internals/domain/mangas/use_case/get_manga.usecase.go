package usecase

import (
	"fmt"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type GetMangaBySlugService struct {
}

type GetMangaBySlugRequest struct {
	Slug string `json:"slug"`
}
type GetMangaBySlugResponse struct {
	Manga entities.Manga
}

func (self *GetMangaBySlugService) Exec(request GetMangaBySlugRequest) (*GetMangaBySlugResponse, error) {
	envs := utils.Utils{}.Envs()
	_ = fmt.Sprintf("%v/%v", envs["base_url"], request.Slug)

	response := GetMangaBySlugResponse{}

	return &response, nil
}
