package events

import "qaanii/shared/entities"

const SEARCH_MANGA_EVENT Events = "@search/query"
const SEARCHED_MANGA_EVENT Events = "@search/results"

type SearchMangaMessage struct {
	Query string `json:"query"`

	BaseEvent
}

type SearchedMangaMessage struct {
	Data []entities.Manga `json:"data"`

	BaseEvent
}
