package events

import "qaanii/shared/entities"

const SCRAPE_CHAPTER_EVENT Events = "@chapter/scrape"
const SCRAPED_CHAPTER_EVENT Events = "@chapter/scraped"

type ScrapeChapterMessage struct {
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`

	BaseEvent
}

type ScrapedChapterMessage struct {
	Data entities.Chapter `json:"data"`

	BaseEvent
}
