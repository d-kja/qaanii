package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"server/internals/domain/scraping/entities"
	"server/internals/http/middleware"
	"server/internals/utils"

	"github.com/gofiber/fiber/v2"
)

type GetChapterService struct {
	Scraper middleware.RodMiddleware
}

type GetChapterRequest struct {
	Slug string
	Ctx  *fiber.Ctx
}

type GetChapterResponse struct {
	Pages []entities.Page `json:"pages"`
}

func (self GetChapterService) Exec(request GetChapterRequest) (*GetChapterResponse, error) {
	browser, page := self.Scraper.New()
	defer browser.MustClose()

	if page == nil {
		return nil, errors.New("Invalid page")
	}
	defer page.MustClose()

	url := fmt.Sprintf("%v/%v", utils.BASE_URL, request.Slug)

	page.MustNavigate(url)
	page.MustWaitLoad()

	self.Scraper.HandleGuard(page)

	chapter_ch := make(chan bool, 1)
	self.Scraper.QueryRetryX(MANGA_CONTAINER, 5, page, chapter_ch)

	is_chapter_loaded := <-chapter_ch
	if !is_chapter_loaded {
		return nil, errors.New("Manga pages not found")
	}

	page_window, err := page.GetWindow()
	if err == nil {
		page_height := page_window.Height

		if page_height != nil {
			offset_y := float64(*page_height) + 5000.0

			time.Sleep(time.Second * 5)
			utils.LOGGER.INFO.Println("Scrolling to preload...")
			page.Mouse.Scroll(0, offset_y, 5)
			time.Sleep(time.Second * 2)

			utils.LOGGER.INFO.Println("Scrolling to preload...")
			offset_y += 15000.0
			page.Mouse.Scroll(0, offset_y, 5)
			time.Sleep(time.Second * 2)

			utils.LOGGER.INFO.Println("Scrolling to preload...")
			offset_y += 15000.0
			page.Mouse.Scroll(0, offset_y, 5)
			time.Sleep(time.Second * 2)

			utils.LOGGER.INFO.Println("Scrolling to preload...")
			offset_y += 15000.0
			page.Mouse.Scroll(0, offset_y, 5)
			time.Sleep(time.Second * 2)
		}
	}

	pages_container := page.MustElementsX(MANGA_IMAGES)
	pages := []entities.Page{}

	for idx, page := range pages_container {
		manga_page := entities.Page{
			Order: idx,
		}

		if page == nil {
			utils.LOGGER.ERROR.Printf("Page (%v) - Invalid container\n", idx+1)
			continue
		}

		image_url, err := page.Attribute("src")
		if err != nil {
			utils.LOGGER.ERROR.Printf("Page (%v) - Invalid url\n", idx+1)
			continue
		}

		if image_url == nil {
			utils.LOGGER.ERROR.Printf("Page (%v) - Image url invalid pointer: %+v\n", idx+1, image_url)
			continue
		}

		manga_page.ImageUrl = *image_url

		image_resource, err := page.Resource()
		if err != nil || len(image_resource) == 0 {
			utils.LOGGER.WARN.Printf("Page (%v) - Unable to retrieve image resource, error(%v): %+v\n", idx+1, len(image_resource), err)
			continue
		}

		manga_page.ImageType = http.DetectContentType(image_resource)
		manga_page.Image = base64.StdEncoding.EncodeToString(image_resource)

		pages = append(pages, manga_page)
	}

	response := GetChapterResponse{
		Pages: pages,
	}

	return &response, nil
}

var (
	MANGA_CONTAINER = "//div[contains(@id, 'chapter-images')]"
	MANGA_IMAGES    = "//div[contains(@id, 'chapter-images')]/div[contains(@class, 'chapter-image')]/img"
)
