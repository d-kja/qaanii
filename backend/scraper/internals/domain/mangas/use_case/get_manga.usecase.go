package usecase

import (
	"errors"
	"fmt"
	"log"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	"qaanii/scraper/internals/domain/mangas/constants"
	"qaanii/shared/entities"
	"qaanii/shared/utils"
)

type GetMangaBySlugService struct {
	Scraper coreentities.Scraper
}

type GetMangaBySlugRequest struct {
	Slug string `json:"slug"`
}
type GetMangaBySlugResponse struct {
	Manga entities.Manga
}

func (self *GetMangaBySlugService) Exec(request GetMangaBySlugRequest) (*GetMangaBySlugResponse, error) {
	envs := utils.Envs()

	url := fmt.Sprintf("%v/%v", envs["base_url"], request.Slug)
	page := self.Scraper.NewPage(url)
	defer page.MustClose()

	was_found := self.Scraper.CheckWithRetry(page, constants.CHAPTER_LINKS, coreentities.RetryConfig{RetryType: coreentities.RETRY_XPATH_MANY, MaxRetries: 5})
	if !was_found {
		log.Println("Manga | Chapter links not found.")
		return nil, errors.New("Chapters not found.")
	}

	el_chapters := page.MustElementsX(constants.CHAPTER_LINKS)
	chapters := []entities.Chapter{}

	for idx, el_chapter := range el_chapters {
		if el_chapter == nil {
			log.Printf("Manga | Chapter [%v] is nil\n", idx)
			continue
		}

		link, err := el_chapter.Attribute("href")
		if err != nil {
			log.Printf("Manga | Chapter [%v] has an invalid url\n", idx)
			continue
		}

		url := (*link)[1:]
		chapter := entities.Chapter{
			Link: url,
		}

		title_container, err := el_chapter.ElementX(constants.CHAPTER_TITLE)
		if err != nil {
			chapter.Title = fmt.Sprintf("Chapter [%v]?", idx+1)
			chapters = append(chapters, chapter)

			log.Printf("Manga | Chapter [%v] has an invalid title\n", idx)
			continue
		}

		title, err := title_container.Text()
		if err != nil {
			chapter.Title = fmt.Sprintf("Chapter [%v]?", idx+1)
			chapters = append(chapters, chapter)

			log.Printf("Manga | Chapter [%v] has an invalid title\n", idx)
			continue
		}

		chapter.Title = title

		time_container, err := el_chapter.ElementX(constants.CHAPTER_DATE)
		if err != nil {
			chapters = append(chapters, chapter)

			log.Printf("Manga | Chapter [%v] has an invalid date/time\n", idx)
			continue
		}

		time, err := time_container.Text()
		if err != nil {
			chapters = append(chapters, chapter)

			log.Printf("Manga | Chapter [%v] has an invalid date/time\n", idx)
			continue
		}

		chapter.Time = time
		chapters = append(chapters, chapter)
	}

	response := GetMangaBySlugResponse{
		Manga: entities.Manga{
			Chapters: &chapters,
		},
	}

	has_status := self.Scraper.CheckWithRetry(page, constants.MANGA_STATUS, coreentities.RetryConfig{RetryType: coreentities.RETRY_XPATH, MaxRetries: 2})
	has_last_updated := self.Scraper.CheckWithRetry(page, constants.MANGA_LAST_UPDATE, coreentities.RetryConfig{RetryType: coreentities.RETRY_XPATH, MaxRetries: 2})

	if has_status {
		manga_status := page.MustElementX(constants.MANGA_STATUS)
		status, err := manga_status.Text()

		if err == nil {
			response.Manga.Status = &status
		}
	}

	if has_last_updated {
		manga_update := page.MustElementX(constants.MANGA_LAST_UPDATE)
		update, err := manga_update.Text()

		if err == nil {
			response.Manga.Time = &update
		}
	}

	return &response, nil
}
