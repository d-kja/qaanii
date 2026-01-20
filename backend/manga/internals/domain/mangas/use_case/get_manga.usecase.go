package usecase

import (
	"qaanii/shared/entities"
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
	// envs := utils.Utils{}.Envs()
	// url := fmt.Sprintf("%v/%v", envs["base_url"], request.Slug)

	// SETUP GRPC TO RETURN A STREAM RESPONSE...

	response := GetMangaBySlugResponse{}

	return &response, nil
}
