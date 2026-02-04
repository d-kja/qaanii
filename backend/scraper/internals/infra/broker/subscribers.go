package broker

import (
	"qaanii/scraper/internals/infra/broker/manga"
	"qaanii/scraper/internals/infra/broker/search"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"
	"sync"
)

func SetupSubscribers(request broker.SubscriberRequest) {
	// Not necessary, but to avoid sync issues during hot-reload I added a wg
	wg := sync.WaitGroup{}
	wg.Add(3)

	go broker.CreateConsumer(events.SEARCH_MANGA_EVENT, request, search.SearchByNameSubscriber, &wg)
	go broker.CreateConsumer(events.SCRAPE_MANGA_EVENT, request, manga.ScrapeMangaSubscriber, &wg)
	go broker.CreateConsumer(events.SCRAPE_CHAPTER_EVENT, request, manga.ScrapeChapterSubscriber, &wg)

	wg.Wait()
}
