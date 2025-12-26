package constants

var (
	LIST_CONTAINER string = "div.list.manga-list"
	MANGA_LIST     string = "//div[contains(@class, 'list') and contains(@class, 'manga-list')]//div[contains(@class, 'book-detailed-item')]"

	MANGA_THUMBNAIL string = ".//div[contains(@class, 'thumb')]//img"
	MANGA_CONTENT   string = ".//div[contains(@class, 'meta')]"

	MANGA_TITLE       string = ".//div[contains(@class, 'title')]/h3/a"
	MANGA_DESCRIPTION string = ".//div[contains(@class, 'summary')]/p"
	MANGA_TAGS        string = ".//div[contains(@class, 'genres')]/span"
)
