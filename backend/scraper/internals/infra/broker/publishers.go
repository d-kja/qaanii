package broker

import (
	"qaanii/shared/broker"
	"qaanii/shared/broker/events"
)

func SetupPublishers(request broker.PublisherRequest) {
	broker.CreatePublisher(events.SEARCHED_MANGA_EVENT, request)
	broker.CreatePublisher(events.SCRAPED_MANGA_EVENT, request)
	broker.CreatePublisher(events.SCRAPED_CHAPTER_EVENT, request)
}
