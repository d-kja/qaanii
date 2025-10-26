package services

import (
	"errors"
	"fmt"

	"server/internals/domain/scraping/entities"
	"server/internals/http/middleware"
	"server/internals/utils"

	"github.com/gofiber/fiber/v2"
)

type GetBySlugService struct {
	Scraper middleware.RodMiddleware
}

type GetBySlugRequest struct {
	Slug string
	Ctx  *fiber.Ctx
}

type GetBySlugResponse struct {
	Chapters []entities.Chapter `json:"chapters"`
}

func (self GetBySlugService) Exec(request GetBySlugRequest) (*GetBySlugResponse, error) {
	scraper, page := self.Scraper.New()
	defer page.MustClose()
	defer scraper.MustClose()

	if page == nil {
		return nil, errors.New("Invalid page")
	}

	url := fmt.Sprintf("%v/%v", utils.BASE_URL, request.Slug)

	page.MustNavigate(url)
	page.MustWaitLoad()

	self.Scraper.HandleGuard(page)

	manga_ch := make(chan bool, 1)
	self.Scraper.QueryManyRetryX(CHAPTER_LINKS, 5, page, manga_ch)

	is_manga_list_available := <-manga_ch
	if !is_manga_list_available {
		return nil, errors.New("Manga list not available")
	}

	chapters_container := page.MustElementsX(CHAPTER_LINKS)
	chapters := []entities.Chapter{}
	for idx, chapter := range chapters_container {
		manga_chapter := entities.Chapter{
			Title: "Not found",
			Date:  "-",
		}

		if chapter == nil {
			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid container\n", idx+1)
			continue
		}

		chapter_url, err := chapter.Attribute("href")
		if err != nil {
			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid url\n", idx+1)
			continue
		}

		manga_chapter.Slug = (*chapter_url)[1:]

		chapter_title_container, err := chapter.ElementX(CHAPTER_TITLE)
		if err != nil {
			chapters = append(chapters, manga_chapter)

			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid title\n", idx+1)
			continue
		}

		title, err := chapter_title_container.Text()
		if err != nil {
			chapters = append(chapters, manga_chapter)

			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid title text\n", idx+1)
			continue
		}

		manga_chapter.Title = title

		chapter_time_container, err := chapter.ElementX(CHAPTER_DATE)
		if err != nil {
			chapters = append(chapters, manga_chapter)

			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid time\n", idx+1)
			continue
		}

		chapter_time, err := chapter_time_container.Text()
		if err != nil {
			chapters = append(chapters, manga_chapter)

			utils.LOGGER.ERROR.Printf("Chapter (%v) - Invalid time text\n", idx+1)
			continue
		}

		manga_chapter.Date = chapter_time
		chapters = append(chapters, manga_chapter)
	}

	response := GetBySlugResponse{
		Chapters: chapters,
	}

	return &response, nil
}

var (
	CHAPTER_LINKS string = "//ul[@id = 'chapter-list']/li/a"
	CHAPTER_TITLE string = "./div/strong[contains(@class, 'chapter-title')]"
	CHAPTER_DATE  string = "./div/time[contains(@class, 'chapter-update')]"
)
