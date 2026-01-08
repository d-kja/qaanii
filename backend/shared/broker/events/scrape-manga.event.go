package events

import "qaanii/shared/entities"

const SCRAPE_MANGA_EVENT Events = "@manga/scrape"
const SCRAPED_MANGA_EVENT Events = "@manga/scraped"

type ScrapeMangaMessage struct {
	Slug string `json:"slug"`

	BaseEvent
}

type ScrapedMangaMessage struct {
	Data entities.Manga `json:"data"`

	BaseEvent
}
