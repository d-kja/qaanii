package coreutils

import (
	"context"
	"errors"
	"log"
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	"sync/atomic"
	"time"
)

const MAX_SCRAPER_INSTANCES uint8 = 3 // Increase this according to your needs.
const DEFAULT_ACQUIRE_TIMEOUT = time.Second * 45

var BROWSER_POOL_EXHAUSTED = errors.New("Unable to retrieve browser instance, pool exhausted.")

type BrowserPool struct {
	closed *atomic.Bool

	Size     uint8
	browsers chan *coreentities.Scraper
}

func NewBrowserPool() BrowserPool {
	pool := BrowserPool{}

	pool.closed = &atomic.Bool{}
	pool.Size = MAX_SCRAPER_INSTANCES
	pool.browsers = make(chan *coreentities.Scraper, MAX_SCRAPER_INSTANCES)

	for range MAX_SCRAPER_INSTANCES {
		scraper := coreentities.NewScraper()
		pool.browsers <- &scraper // Pre allocate browsers
	}

	return pool
}

func (self *BrowserPool) Get() (*coreentities.Scraper, error) {
	timeout_ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_ACQUIRE_TIMEOUT)
	defer cancel()

	select {
	case browser := <-self.browsers:
		{
			return browser, nil
		}
	case <-timeout_ctx.Done():
		{
			log.Println("[BROWSER_POOL] - Resource timed out, unable to retrieve a browser")
			return nil, BROWSER_POOL_EXHAUSTED
		}
	}
}

func close_pages(scraper *coreentities.Scraper) {
	for _, page := range scraper.Browser.MustPages() {
		page.MustClose()
	}
}

func (self *BrowserPool) Release(scraper *coreentities.Scraper) {
	if self.closed.Load() {
		scraper.Browser.MustClose()
		return
	}

	close_pages(scraper)

	select {
	case self.browsers <- scraper:
		{
			return
		}

	// Queue full (> MAX_SCRAPER_INSTANCES)
	default:
		{
			// Clean up
			scraper.Browser.MustClose()
		}
	}
}

// Free resources
func (self *BrowserPool) Cleanup() {
	if self.closed.Swap(true) {
		return
	}

	close(self.browsers)

	for scraper := range self.browsers {
		close_pages(scraper)
		scraper.Browser.MustClose()
	}
}
