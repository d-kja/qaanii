package entities

type Page struct {
	Order     int    `json:"order"`
	Image     string `json:"image"`
	ImageUrl  string `json:"image_url"`
	ImageType string `json:"image_type"`
}
