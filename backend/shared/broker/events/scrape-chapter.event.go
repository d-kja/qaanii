package events

import "qaanii/shared/entities"

const SCRAPE_CHAPTER_EVENT Events = "@chapter/scrape"
const SCRAPED_CHAPTER_EVENT Events = "@chapter/scraped"

type ScrapeChapterMessage struct {
	Slug    string `json:"slug"`
	Chapter string `json:"chapter"`

	BaseEvent
}

type MessageChapter struct {
	Title string `json:"title"`
	Link  string `json:"link"`

	Time string `json:"time"`

	Pages []entities.Page `json:"pages"`
}

type ScrapedChapterMessage struct {
	Data MessageChapter `json:"data"`

	BaseEvent
}
