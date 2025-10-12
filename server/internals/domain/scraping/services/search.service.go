package services

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"server/internals/domain/scraping/entities"
	"server/internals/http/middleware"
	"server/internals/utils"
)

type SearchMangasService struct {
	Scraper middleware.RodMiddleware
}

type SearchMangasRequest struct {
	Query string
}

type SearchMangasResponse struct {
	Mangas []entities.Manga
}

func (self SearchMangasService) Exec(request SearchMangasRequest) (*SearchMangasResponse, error) {
	scraper, page := self.Scraper.New()
	defer scraper.MustClose()

	url := fmt.Sprintf("%v/search?q=%v", utils.BASE_URL, url.QueryEscape(request.Query))

	page.MustNavigate(url)
	page.MustWaitLoad()

	self.Scraper.HandleGuard(page)

	// Wait manga list to load
	page.MustElement(LIST_CONTAINER)

	mangas_element := page.MustElementsX(MANGA_LIST)
	mangas := []entities.Manga{}

	for idx, manga_container := range mangas_element {
		manga := entities.Manga{}

		content, content_err := manga_container.ElementX(MANGA_CONTENT)
		if content_err != nil {
			utils.LOGGER.ERROR.Printf("Manga (%v) - Content not found, error: %+v\n", idx+1, content_err)
			continue
		}

		// Process main content
		manga_header, header_err := content.ElementX(MANGA_TITLE)
		if header_err != nil {
			utils.LOGGER.WARN.Printf("Manga (%v) - Title element not found, error: %+v\n", idx+1, header_err)
			continue
		}

		manga_title, title_err := manga_header.Attribute("title")
		if title_err != nil {
			utils.LOGGER.WARN.Printf("Manga (%v) - Title Attr not found, error: %+v\n", idx+1, title_err)
			continue
		}

		manga_url, url_err := manga_header.Attribute("href")
		if url_err != nil {
			utils.LOGGER.WARN.Printf("Manga (%v) - URL Attr not found, error: %+v\n", idx+1, url_err)
			continue
		}

		manga.Name = *manga_title
		manga.Url = fmt.Sprintf("%v%v", utils.BASE_URL, *manga_url)

		manga_description, description_err := content.ElementX(MANGA_DESCRIPTION)
		if description_err != nil {
			utils.LOGGER.WARN.Printf("Manga (%v) - Description not found, error: %+v\n", idx+1, description_err)
		} else {
			description, err := manga_description.Text()

			if err == nil {
				manga.Description = description
			}
		}

		manga_tags, tags_err := content.ElementsX(MANGA_TAGS)
		if tags_err != nil {
			utils.LOGGER.WARN.Printf("Manga (%v) - Tags not found, error: %+v\n", idx+1, tags_err)
		} else {
			tags := []string{}

			for _, tag := range manga_tags {
				content, err := tag.Text()
				if err != nil {
					continue
				}

				tags = append(tags, content)
			}
		}

		thumbnail_container, thumb_err := manga_container.ElementX(MANGA_THUMBNAIL)
		if thumb_err != nil {
			mangas = append(mangas, manga)

			utils.LOGGER.WARN.Printf("Manga (%v) - Thumbnail element not found, error: %+v\n", idx+1, thumb_err)
			continue
		}

		// Extract image (optional)
		thumbnail_src, err := thumbnail_container.Attribute("src")
		if err != nil || thumbnail_src == nil {
			mangas = append(mangas, manga)

			utils.LOGGER.WARN.Printf("Manga (%v) - Thumbnail src not found, error: %+v\n", idx+1, err)
			continue
		}

		manga.ImageUrl = *thumbnail_src
		
		image_resource, err := thumbnail_container.Resource()
		if err != nil || len(image_resource) == 0 {
			mangas = append(mangas, manga)

			utils.LOGGER.WARN.Printf("Manga (%v) - Unable to retrieve image resource, error(%v): %+v\n", idx+1, len(image_resource), err)
			continue
		}

		manga.ImageType = http.DetectContentType(image_resource)
		manga.Image = base64.StdEncoding.EncodeToString(image_resource)
		mangas = append(mangas, manga)
	}

	response := SearchMangasResponse{
		Mangas: mangas,
	}

	return &response, nil
}

var (
	LIST_CONTAINER string = "div.list.manga-list"
	MANGA_LIST     string = "//div[contains(@class, 'list') and contains(@class, 'manga-list')]//div[contains(@class, 'book-detailed-item')]"

	MANGA_THUMBNAIL string = ".//div[contains(@class, 'thumb')]//img"
	MANGA_CONTENT   string = ".//div[contains(@class, 'meta')]"

	MANGA_TITLE       string = ".//div[contains(@class, 'title')]/h3/a"
	MANGA_DESCRIPTION string = ".//div[contains(@class, 'summary')]/p"
	MANGA_TAGS        string = ".//div[contains(@class, 'genres')]/span"
)
