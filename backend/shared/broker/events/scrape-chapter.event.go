package events

import "qaanii/shared/entities"

const SCRAPE_CHAPTER_EVENT Events = "@manga/scrape-chapter"
const SCRAPED_CHAPTER_EVENT Events = "@manga/scraped-chapter"

type ScrapeChapterMessage struct {
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`

	BaseEvent
}

type ScrapedChapterMessage struct {
	Data entities.Chapter `json:"data"`

	BaseEvent
}
