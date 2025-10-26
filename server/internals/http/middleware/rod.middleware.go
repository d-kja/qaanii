package middleware

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"server/internals/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

const ROD_KEY string = "@locals/rod"

var ROD_AGENTS []string = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0",
	"Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0",
	"Mozilla/5.0 (Linux; Android 14; 24030PN60G Build/UKQ1.231003.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/122.0.6261.119 Mobile Safari/537.36",
}

type RodMiddleware struct{}

func (self RodMiddleware) New() (*rod.Browser, *rod.Page) {
	scraper := rod.
		New().
		Timeout(time.Minute).
		MustConnect()

	utils.LOGGER.INFO.Println("Setting up Stealth")
	page := self.SetupStealth(scraper)

	return scraper, page
}

func (RodMiddleware) SetupStealth(instance *rod.Browser) *rod.Page {
	page := stealth.MustPage(instance)
	return page
}

func (RodMiddleware) QueryRetry(query string, max_retries int, page *rod.Page, query_ch chan bool) {
	retry_count := 0
	found := false

	for {
		utils.LOGGER.INFO.Printf("[RETRY/SINGLE] - Found: %v\n", found)
		if found {
			break
		}

		if retry_count >= max_retries {
			break
		}

		ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
		defer cancel()

		thread_ch := make(chan bool, 1)
		defer close(thread_ch)

		go func(ch chan bool, ctx context.Context) {
			element, err := page.Element(query)

			canceled_ctx := func(ctx context.Context) bool {
				if err := ctx.Err(); err != nil {
					utils.LOGGER.WARN.Println("Context was cancelled, returning")
					return true
				}

				return false
			}

			if err != nil {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			is_visible, err := element.Visible()
			if err != nil {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			if !is_visible {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			was_canceled := canceled_ctx(ctx)
			if was_canceled {
				return
			}

			thread_ch <- true
		}(thread_ch, ctx_timeout)

		select {
		case <-ctx_timeout.Done():
			{
				utils.LOGGER.INFO.Printf("Context timed out, retrying... (Count: %v) \n", retry_count)
				retry_count++

				continue
			}
		case was_found := <-thread_ch:
			{
				found = was_found
			}
		}
	}

	query_ch <- found
}

func (RodMiddleware) QueryRetryX(query string, max_retries int, page *rod.Page, query_ch chan bool) {
	retry_count := 0
	found := false

	for {
		utils.LOGGER.INFO.Printf("[RETRY/SINGLEX] - Found: %v\n", found)
		if found {
			break
		}

		if retry_count >= max_retries {
			utils.LOGGER.INFO.Println("[RETRY/SINGLEX] - Max retry count reached, stopping...")
			break
		}

		ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
		defer cancel()

		thread_ch := make(chan bool, 1)
		defer close(thread_ch)

		go func(ch chan bool, ctx context.Context) {
			element, err := page.ElementX(query)

			canceled_ctx := func(ctx context.Context) bool {
				if err := ctx.Err(); err != nil {
					utils.LOGGER.WARN.Println("Context was cancelled, returning...")
					return true
				}

				return false
			}

			if err != nil {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				utils.LOGGER.INFO.Printf("[RETRY/SINGEX] - An error occurred while querying the element, error: %v\n", err)
				thread_ch <- false
				return
			}

			text, err := element.Text()
			if err != nil {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			if len(text) == 0 {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			was_canceled := canceled_ctx(ctx)
			if was_canceled {
				utils.LOGGER.INFO.Println("[RETRY/SINGLEX] - Element found, but the context was canceled")
				return
			}

			thread_ch <- true
		}(thread_ch, ctx_timeout)

		select {
		case <-ctx_timeout.Done():
			{
				utils.LOGGER.INFO.Printf("[RETRY/SINGLEX] - Context timed out, retrying... (Count: %v) \n", retry_count + 1)
				retry_count++

				continue
			}
		case was_found := <-thread_ch:
			{
				utils.LOGGER.INFO.Printf("[RETRY/SINGLEX] - Thread returned a response, found: %v\n", was_found)
				found = was_found
			}
		}
	}

	query_ch <- found
}

func (RodMiddleware) QueryManyRetryX(query string, max_retries int, page *rod.Page, query_ch chan bool) {
	retry_count := 0
	found := false

	for {
		utils.LOGGER.INFO.Printf("[RETRY/MANYX] - Found: %v\n", found)
		if found {
			break
		}

		if retry_count >= max_retries {
			break
		}

		ctx_timeout, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
		defer cancel()

		thread_ch := make(chan bool, 1)
		defer close(thread_ch)

		go func(ch chan bool, ctx context.Context) {
			element, err := page.ElementsX(query)

			canceled_ctx := func(ctx context.Context) bool {
				if err := ctx.Err(); err != nil {
					utils.LOGGER.WARN.Println("Context was cancelled, returning")
					return true
				}

				return false
			}

			if err != nil {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			is_empty := element.Empty()
			if is_empty {
				retry_count++

				was_canceled := canceled_ctx(ctx)
				if was_canceled {
					return
				}

				thread_ch <- false
				return
			}

			was_canceled := canceled_ctx(ctx)
			if was_canceled {
				return
			}

			thread_ch <- true
		}(thread_ch, ctx_timeout)

		select {
		case <-ctx_timeout.Done():
			{
				utils.LOGGER.INFO.Printf("Context timed out, retrying... (Count: %v) \n", retry_count)
				retry_count++

				continue
			}
		case was_found := <-thread_ch:
			{
				if !was_found {
					time.Sleep(time.Second * 2)
				}

				found = was_found
			}
		}
	}

	query_ch <- found
}

func (self RodMiddleware) HandleGuard(page *rod.Page) {
	query := "//title[contains(translate(text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'ddos-guard')]"

	query_ch := make(chan bool, 1)
	defer close(query_ch)

	utils.LOGGER.INFO.Printf("[RETRY] - Searching for Guard\n")
	self.QueryRetryX(query, 5, page, query_ch)
	was_found := <-query_ch

	if !was_found {
		utils.LOGGER.INFO.Printf("[RETRY] - Guard not found [Safe query]\n")
		return
	}

	element, err := page.ElementX(query)
	if err != nil {
		utils.LOGGER.INFO.Printf("[BLOCKING] - Guard not found\n")
		return
	}

	if element == nil {
		utils.LOGGER.INFO.Printf("Guard pointer nil\n")
		return
	}

	content, err := element.Text()
	if err != nil {
		utils.LOGGER.INFO.Printf("Guard found, but content unavailable. %+v\n", err)
		return
	}

	utils.LOGGER.INFO.Printf("Guard found, type: %v. Waiting 5 seconds...\n", content)
	time.Sleep(time.Second * 5)
}

func (RodMiddleware) GetAgent() string {
	agents_len := len(ROD_AGENTS) - 1
	agent_idx := rand.IntN(agents_len)

	return ROD_AGENTS[agent_idx]
}

func (RodMiddleware) Metadata(page *rod.Page) {
	element := page.MustElement("#broken-image-dimensions.passed")

	for _, row := range element.MustParents("table").First().MustElements("tr:nth-child(n+2)") {
		cells := row.MustElements("td")
		key := cells[0].MustProperty("textContent")

		if strings.HasPrefix(key.String(), "User Agent") {
			fmt.Printf("\t\t%s: %t\n", key, !strings.Contains(cells[1].MustProperty("textContent").String(), "HeadlessChrome/"))
		} else if strings.HasPrefix(key.String(), "Hairline Feature") {
			// Detects support for hidpi/retina hairlines, which are CSS borders with less than 1px in width, for being physically 1px on hidpi screens.
			// Not all the machine suppports it.
			continue
		} else {
			fmt.Printf("\t\t%s: %s\n", key, cells[1].MustProperty("textContent"))
		}
	}

	page.MustScreenshot("")
}
