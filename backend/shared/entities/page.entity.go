package entities

type Page struct {
	Order     int    `json:"order"`

	Image     string `json:"image"`
	ImageType string `json:"image_type"`
}
