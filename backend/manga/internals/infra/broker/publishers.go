package broker

import (
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"
)

func SetupPublishers(request broker.PublisherRequest) {
	broker.CreatePublisher(events.SCRAPE_CHAPTER_EVENT, request)
	broker.CreatePublisher(events.SCRAPE_MANGA_EVENT, request)
	broker.CreatePublisher(events.SEARCH_MANGA_EVENT, request)
}
