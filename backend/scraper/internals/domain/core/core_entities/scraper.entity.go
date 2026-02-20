package coreentities

import (
	"context"
	"log"
	"math"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Scraper struct {
	Browser *rod.Browser
	Page    *rod.Page
}

func NewScraper() Scraper {
	extension_path, err := filepath.Abs("fixtures/chrome-extensions/ublock")
	if err != nil {
		log.Panicf("[SCRAPER] - Unable to load extensions, error: %+v\n", err)
	}

	instance := launcher.New().
		Set("load-extension", extension_path).
		Headless(false).
		MustLaunch()

	scraper := rod.New().ControlURL(instance).MustConnect().NoDefaultDevice()

	return Scraper{
		Browser: scraper,
	}
}

func (self *Scraper) NewPage(url string) *rod.Page {
	log.Printf("[SCRAPPER] - Creating new page, base url [%v]\n", url)

	page := stealth.MustPage(self.Browser).MustWaitStable()
	self.Page = page.MustNavigate(url).MustWaitStable()

	self.HandleGuard(page)
	self.HandleModals(page)

	return self.Page
}

func (self *Scraper) HandleGuard(page *rod.Page) {
	log.Println("[SCRAPPER] - Searching for Guard Elements")

	was_found := self.CheckWithRetry(page, GUARD_QUERY, RetryConfig{RetryType: RETRY_XPATH, MaxRetries: 2})
	if !was_found {
		log.Println("[SCRAPPER] - Not found, skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Element found, validating query")
	element, err := page.ElementX(GUARD_QUERY)
	if err != nil || element == nil {
		log.Println("[SCRAPPER] - Guard not found, skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Query validated, testing content")
	content, err := element.Text()
	if err != nil || len(content) == 0 {
		log.Println("[SCRAPPER] - Guard found, but content empty skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Content tested, waiting 5 seconds for the guard")
	time.Sleep(time.Second * 5)
}

func (self *Scraper) HandleModals(page *rod.Page) {
	log.Println("[SCRAPPER] - Searching for WARNING modal")

	was_found := self.CheckWithRetry(page, WARNING_QUERY, RetryConfig{RetryType: RETRY_XPATH, MaxRetries: 2})
	if !was_found {
		log.Println("[SCRAPPER] - Not found, skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Element found, validating query")
	element, err := page.ElementX(WARNING_QUERY)
	if err != nil || element == nil {
		log.Println("[SCRAPPER] - Warning modal not found, skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Query validated, accepting warning")
	content, err := element.Text()
	if err != nil || len(content) == 0 {
		log.Println("[SCRAPPER] - Warning modal found, but content empty skipping logic")
		return
	}

	btn_element, err := page.ElementX(WARNING_CONFIRM_BTN_QUERY)
	if err != nil {
		log.Println("[SCRAPPER] - Warning modal confirm button not found, skipping logic")
		return
	}

	err = btn_element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		log.Println("[SCRAPPER] - An error occurred while clicking the warning modal confirm button, skipping logic")
		return
	}

	log.Println("[SCRAPPER] - Modal closed successfully, waiting 3 seconds")
	time.Sleep(time.Second * 3)
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
					log.Println("[SCRAPPER] - Manga Chapter | Scroll behavior - reached scroll limit")
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
			log.Printf("[SCRAPPER] - Queried item not found.\n")
			break
		}

		if was_found {
			log.Printf("[SCRAPPER] - Queried item found!\n")
			break
		}

		func() {
			found_ch := make(chan bool, 1)
			ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Second*2)

			// Clean up
			defer func() {
				str_len := len(query)
				idx := int(math.Min(float64(str_len), 40))

				log.Printf("[SCRAPPER] - Cleaning up, base query [%v...] [%v]\n", query[:idx], count)

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
					log.Printf("[SCRAPPER] - Retry channel finished, was element found? [%v]\n", was_found)
					break
				}

			case <-ctx_timeout.Done():
				{
					log.Printf("[SCRAPPER] - Context timed out, count [%v]\n", count)
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
const WARNING_QUERY string = "//h5[contains(translate(text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'adult content warning')]"
const WARNING_CONFIRM_BTN_QUERY string = "//h5[contains(translate(text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'adult content warning')]/../../div[contains(@class, 'modal-footer')]/button[contains(@class, 'btn-warning')]"
