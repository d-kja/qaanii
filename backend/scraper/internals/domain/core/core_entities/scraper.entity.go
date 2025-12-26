package coreentities

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

type Scraper struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func NewScraper() Scraper {
	scraper := rod.New().MustConnect().NoDefaultDevice()

	return Scraper{
		Browser: scraper,
	}
}

func (self *Scraper) NewPage(url string) *rod.Page {
	log.Printf("Creating new page, base url [%v]\n", url)

	page := stealth.MustPage(self.Browser).MustWaitStable()
	self.Page = page.MustNavigate(url).MustWaitStable()

	self.HandleGuard(page)
	return self.Page
}

func (self *Scraper) HandleGuard(page *rod.Page) {
	log.Println("Searching for Guard Elements")

	was_found := self.CheckWithRetry(page, GUARD_QUERY, RetryConfig{RetryType: RETRY_XPATH, MaxRetries: 2})
	if !was_found {
		log.Println("Not found, skipping logic")
		return
	}

	log.Println("Element found, validating query")
	element, err := page.ElementX(GUARD_QUERY)
	if err != nil || element == nil {
		log.Println("Guard not found, skipping logic")
		return
	}

	log.Println("Query validated, testint content")
	content, err := element.Text()
	if err != nil || len(content) == 0 {
		log.Println("Guard found, but content empty skipping logic")
		return
	}

	log.Println("Content tested, waiting 5 seconds for the guard")
	time.Sleep(time.Second * 5)
}

type PreloadWithScrollConfig struct {
	ScrollLimit   int
	ScrollSpeed   int
	ScrollTimeout int
}

func (Scraper) PreloadWithScroll(page *rod.Page, config PreloadWithScrollConfig) {
	page_window, err := page.GetWindow()

	if err == nil {
		scroll := func(page *rod.Page, offset_y float64) {
			page.Mouse.Scroll(0, offset_y, config.ScrollSpeed)
			time.Sleep(time.Millisecond * time.Duration(config.ScrollTimeout))
		}

		if page_window.Height != nil {
			last_y := -1

			for range config.ScrollLimit {
				page_window, err = page.GetWindow()
				if err != nil {
					continue
				}

				page_height := *page_window.Height

				page_y := page.MustEval(`() => window.scrollY`).Int()
				if page_y == last_y {
					log.Println("Manga Chapter | Scroll behavior - reached scroll limit")
					break
				}

				// Preload content slowly
				scroll(page, float64(page_height))
				last_y = page_y
			}
		}
	}
}

type RetryConfig struct {
	RetryType  RetryType
	MaxRetries int
}

func (Scraper) CheckWithRetry(page *rod.Page, query string, config RetryConfig) bool {
	was_found := false
	count := 0

	for {
		if count >= config.MaxRetries {
			log.Printf("Queried item not found.\n")
			break
		}

		if was_found {
			log.Printf("Queried item found!\n")
			break
		}

		func() {
			found_ch := make(chan bool, 1)
			ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Second*2)

			// Clean up
			defer func() {
				str_len := len(query)
				idx := int(math.Min(float64(str_len), 40))

				log.Printf("Cleaning up, base query [%v...] [%v]\n", query[:idx], count)

				close(found_ch)
				cancel()

				count++
			}()

			go func() {
				switch config.RetryType {
				case RETRY:
					{
						_, err := page.Element(query)
						if err != nil {
							break
						}

						found_ch <- true
						break
					}

				case RETRY_MANY:
					{
						elements, err := page.Elements(query)
						if err != nil {
							break
						}

						if elements.Empty() {
							break
						}

						found_ch <- true
						break
					}

				case RETRY_XPATH:
					{
						_, err := page.ElementX(query)
						if err != nil {
							break
						}

						found_ch <- true
						break
					}

				case RETRY_XPATH_MANY:
					{
						elements, err := page.ElementsX(query)
						if err != nil {
							break
						}

						if elements.Empty() {
							break
						}

						found_ch <- true
						break
					}
				}
			}()

			select {
			case was_found = <-found_ch:
				{
					log.Printf("Retry channel finished, was element found? [%v]\n", was_found)
					break
				}

			case <-ctx_timeout.Done():
				{
					log.Printf("Context timed out, count [%v]\n", count)
					break
				}
			}
		}()
	}

	return was_found
}

type RetryType = string

var (
	RETRY            RetryType = "RETRY"
	RETRY_MANY       RetryType = "RETRY_MANY"
	RETRY_XPATH      RetryType = "RETRY_XPATH"
	RETRY_XPATH_MANY RetryType = "RETRY_XPATH_MANY"
)

const ROD_KEY string = "@locals/rod"
const GUARD_QUERY string = "//title[contains(translate(text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'ddos-guard')]"
