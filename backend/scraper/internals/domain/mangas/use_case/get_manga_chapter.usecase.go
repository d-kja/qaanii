package usecase

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	"qaanii/scraper/internals/domain/mangas/constants"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type GetMangaChapterService struct {
	Scraper coreentities.Scraper
}

type GetMangaChapterRequest struct {
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`
}
type GetMangaChapterResponse struct {
	Chapter entities.Chapter
}

func (self *GetMangaChapterService) Exec(request GetMangaChapterRequest) (*GetMangaChapterResponse, error) {
	defer self.Scraper.Browser.MustClose()
	envs := utils.Envs()

	url := fmt.Sprintf("%v/%v/%v", envs["base_url"], request.Slug, request.Chapter)
	page := self.Scraper.NewPage(url)
	defer page.MustClose()

	was_found := self.Scraper.CheckWithRetry(page, constants.MANGA_CHAPTER_IMAGES, coreentities.RetryConfig{RetryType: coreentities.RETRY_XPATH_MANY, MaxRetries: 5})
	if !was_found {
		log.Println("Manga Chapter | Container not found.")
		return nil, errors.New("Chapter container not found.")
	}

	self.Scraper.PreloadWithScroll(page, coreentities.PreloadWithScrollConfig{
		ScrollTimeout: 100,
		ScrollSpeed:   10,
		ScrollLimit:   50,
	})

	pages_container := page.MustElementsX(constants.MANGA_CHAPTER_IMAGES)
	pages := []entities.Page{}

	for idx, el_page := range pages_container {
		if el_page == nil {
			log.Printf("Chapter | Page [%v] is nil\n", idx)
			continue
		}

		image_url, err := el_page.Attribute("src")
		if err != nil {
			log.Printf("Chapter | Page [%v] has an invalid url\n", idx+1)
			continue
		}

		if image_url == nil {
			log.Printf("Chapter | Page [%v] image url invalid pointer: %+v\n", idx+1, image_url)
			continue
		}

		image_resource, err := el_page.Resource()
		if err != nil || len(image_resource) == 0 {
			log.Printf("Chapter | Page [%v] unable to retrieve image resource, error [%v]: %+v\n", idx+1, len(image_resource), err)
			continue
		}

		image_type := http.DetectContentType(image_resource)
		image := base64.StdEncoding.EncodeToString(image_resource)

		page := entities.Page{
			Image:     image,
			ImageType: image_type,
		}

		pages = append(pages, page)
	}

	response := GetMangaChapterResponse{
		Chapter: entities.Chapter{
			Title: request.Chapter,
			Pages: &pages,
		},
	}

	return &response, nil
}
