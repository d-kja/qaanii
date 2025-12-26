package entities

type Manga struct {
	Url  string `json:"url"` // INFO: Original path (useful for mapping chapters)
	Slug string `json:"slug"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Tags []string `json:"tags"`

	Image     string `json:"image"` // INFO: Base 64 - I will preload and persist locally, and I don't want to use a blob storage
	ImageType string `json:"image_type"`

	Status *string `json:"status"`
	Time   *string `json:"last_update"`

	Chapters *[]Chapter `json:"chapters"`
}
