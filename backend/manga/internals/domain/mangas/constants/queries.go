package constants

var (
	MANGA_STATUS      string = "//div[contains(@class, 'detail')]/div[contains(@class, 'meta')]/p/strong[contains(text(), 'Status')]/../a/span"
	MANGA_LAST_UPDATE string = "//div[contains(@class, 'detail')]/div[contains(@class, 'meta')]/p/strong[contains(text(), 'Last update')]/../span"

	CHAPTER_LINKS string = "//ul[@id = 'chapter-list']/li/a"
	CHAPTER_TITLE string = "./div/strong[contains(@class, 'chapter-title')]"
	CHAPTER_DATE  string = "./div/time[contains(@class, 'chapter-update')]"
)

var (
	MANGA_CHAPTER_CONTAINER = "//div[contains(@id, 'chapter-images')]"
	MANGA_CHAPTER_IMAGES    = "//div[contains(@id, 'chapter-images')]/div[contains(@class, 'chapter-image')]/img"
)
