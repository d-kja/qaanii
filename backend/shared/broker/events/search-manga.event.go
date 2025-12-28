package events

import "qaanii/shared/entities"

const SEARCH_MANGA_EVENT Events = "@manga/search-manga"
const SEARCHED_MANGA_EVENT Events = "@manga/searched-manga"

type SearchMangaMessage struct {
	Query string `json:"query"`

	BaseEvent
}

type SearchedMangaMessage struct {
	Data []entities.Manga `json:"data"`

	BaseEvent
}
