package services

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	defer scraper.MustClose()

	if page == nil {
		return nil, errors.New("Invalid page")
	}

	url := fmt.Sprintf("%v/%v", utils.BASE_URL, request.Slug)

	page.MustNavigate(url)
	page.MustWaitLoad()

	self.Scraper.HandleGuard(page)

	time.Sleep(time.Second * 5)
	// Wait manga list to load
	count := 0

	for {
		ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		defer cancel()

		found_ch := make(chan bool, 1)
		defer close(found_ch)

		go func() {
			utils.LOGGER.INFO.Println("Waiting for page element")
			elements := page.MustElementsX(CHAPTER_LINKS)

			if !elements.Empty() {
				found_ch <- false
				return
			}

			found_ch <- true
		}()

		select {
		case <-ctx_timeout.Done():
			{
				utils.LOGGER.INFO.Printf("Element not found, repeating count: %v\n", count+1)
				count++
			}
		}

		was_found := <-found_ch

		if was_found {
			utils.LOGGER.INFO.Println("Element found!")
			break
		}

		if count >= 5 {
			utils.LOGGER.INFO.Println("Context reached limit")
			break
		}
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
