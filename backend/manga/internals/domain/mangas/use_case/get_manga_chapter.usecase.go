package usecase

import (
	"fmt"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type GetMangaChapterService struct {
}

type GetMangaChapterRequest struct {
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`
}
type GetMangaChapterResponse struct {
	Chapter entities.Chapter
}

func (self *GetMangaChapterService) Exec(request GetMangaChapterRequest) (*GetMangaChapterResponse, error) {
	envs := utils.Utils{}.Envs()
	_ = fmt.Sprintf("%v/%v/%v", envs["base_url"], request.Slug, request.Chapter)

	response := GetMangaChapterResponse{}

	return &response, nil
}
