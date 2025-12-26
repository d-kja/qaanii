package entities

type Chapter struct {
	Title string `json:"title"`
	Link  string `json:"link"`

	// Plain string, extracted from the website
	Time string `json:"time"`

	Pages *[]Page `json:"pages"`
}
