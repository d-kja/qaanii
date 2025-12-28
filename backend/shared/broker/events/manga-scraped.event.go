package events

const MANGA_SCRAPED_EVENT Events = "@manga/manga-scraped"

type MangaScrapedMessage struct{
	BaseEvent
}
