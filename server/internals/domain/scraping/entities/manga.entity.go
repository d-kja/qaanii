package entities

type Manga struct {
	Image string `json:"image"`
	ImageUrl string `json:"image_url"`
	ImageType string `json:"image_type"`

	Name string `json:"name"`
	Description string `json:"description"`
	
	Tags []string `json:"tags"`
	Url string `json:"url"`
}
