package events

const SCRAPE_MANGA_EVENT Events = "@manga/scrape-manga"

type ScrapeMangaMessage struct {
	BaseEvent
}
