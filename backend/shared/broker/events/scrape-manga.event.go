package events

import "qaanii/shared/entities"

const SCRAPE_MANGA_EVENT Events = "@manga/scrape-manga"
const SCRAPED_MANGA_EVENT Events = "@manga/scraped-manga"

type ScrapeMangaMessage struct {
	Slug string `json:"slug"`

	BaseEvent
}

type ScrapedMangaMessage struct {
	Data entities.Manga `json:"data"`

	BaseEvent
}
