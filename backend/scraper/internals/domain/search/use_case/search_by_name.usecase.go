package usecase

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	"qaanii/scraper/internals/domain/search/constants"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type SearchByNameService struct {
	Scraper coreentities.Scraper
}

type SearchByNameRequest struct {
	Search string
}
type SearchByNameResponse struct {
	Mangas []entities.Manga
}

func (self *SearchByNameService) Exec(request SearchByNameRequest) (*SearchByNameResponse, error) {
	defer self.Scraper.Browser.MustClose()
	envs := utils.Utils{}.Envs()

	search := url.QueryEscape(request.Search)
	url := fmt.Sprintf("%v/search?q=%v", envs["base_url"], search)

	page := self.Scraper.NewPage(url)
	defer page.MustClose()

	chap_list_exists := self.Scraper.CheckWithRetry(page, constants.LIST_CONTAINER, coreentities.RetryConfig{
		RetryType:  coreentities.RETRY,
		MaxRetries: 5,
	})

	if !chap_list_exists {
		return nil, errors.New("Search | Container list not found.")
	}

	page_window, err := page.GetWindow()
	if err == nil {
		page_height := page_window.Height

		if page_height != nil {
			offset_y := float64(*page_height) + 1000.0

			log.Println("Search | Scrolling container list.")
			page.Mouse.Scroll(0, offset_y, 25)
		}
	}

	manga_elements := page.MustElementsX(constants.MANGA_LIST)
	mangas := []entities.Manga{}

	for idx, manga_element := range manga_elements {
		if manga_element == nil {
			log.Printf("Search | Manga element [%v] is nil\n", idx)
			continue
		}

		content, content_err := manga_element.ElementX(constants.MANGA_CONTENT)
		if content_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid content, error: %+v\n", idx, content_err)
			continue
		}

		manga_header, header_err := content.ElementX(constants.MANGA_TITLE)
		if header_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid header, error: %+v\n", idx, header_err)
			continue
		}

		manga_title, title_err := manga_header.Attribute("title")
		if title_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid title, error: %+v\n", idx, title_err)
			continue
		}

		manga_url, url_err := manga_header.Attribute("href")
		if url_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid url, error: %+v\n", idx, url_err)
			continue
		}

		manga := entities.Manga{
			Name: *manga_title,
			Slug: (*manga_url)[1:],
			Url:  fmt.Sprintf("%v%v", envs["base_url"], *manga_url),
		}

		manga_description, description_err := content.ElementX(constants.MANGA_DESCRIPTION)
		if description_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid description, error: %+v\n", idx, description_err)
		} else {
			description, err := manga_description.Text()

			if err == nil {
				manga.Description = description
			}
		}

		manga_tags, tags_err := content.ElementsX(constants.MANGA_TAGS)
		if tags_err != nil {
			log.Printf("Search | Manga element [%v] has an invalid tag section, error: %+v\n", idx, tags_err)
		} else {
			tags := []string{}

			for _, tag := range manga_tags {
				content, err := tag.Text()
				if err != nil {
					continue
				}

				tags = append(tags, content)
			}

			manga.Tags = tags
		}

		thumbnail_container, thumb_err := manga_element.ElementX(constants.MANGA_THUMBNAIL)
		if thumb_err != nil {
			mangas = append(mangas, manga)

			log.Printf("Search | Manga element [%v] has an invalid thumbnail, error: %+v\n", idx, thumb_err)
			continue
		}

		image_resource, err := thumbnail_container.Resource()
		if err != nil || len(image_resource) == 0 {
			mangas = append(mangas, manga)

			log.Printf("Search | Manga element [%v] has an invalid image resource, error: %+v\n", idx, err)
			continue
		}

		manga.ImageType = http.DetectContentType(image_resource)
		manga.Image = base64.StdEncoding.EncodeToString(image_resource)

		mangas = append(mangas, manga)
	}

	response := SearchByNameResponse{
		Mangas: mangas,
	}

	return &response, nil
}
