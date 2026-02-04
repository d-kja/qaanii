package broker

import (
	"qaanii/scraper/internals/infra/broker/manga"
	"qaanii/scraper/internals/infra/broker/search"
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"
)

func SetupSubscribers(request broker.SubscriberRequest) {
	broker.CreateConsumer(events.SEARCH_MANGA_EVENT, request, search.SearchByNameSubscriber)
	broker.CreateConsumer(events.SCRAPE_MANGA_EVENT, request, manga.ScrapeMangaSubscriber)
	broker.CreateConsumer(events.SCRAPE_CHAPTER_EVENT, request, manga.ScrapeChapterSubscriber)
}
