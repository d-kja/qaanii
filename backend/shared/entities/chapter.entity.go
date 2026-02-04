package entities

import mangav1 "qaanii/mangabuf/gen/manga/v1"

type Chapter struct {
	Title string `json:"title"`
	Link  string `json:"link"`

	// Plain string, extracted from the website
	Time string `json:"time"`

	Pages *[]Page `json:"pages"`
}

func (self Chapter) ToProtobuf() mangav1.Chapter {
	return mangav1.Chapter{
		Pages: nil, // TODO: No one deservers this.
		Title: self.Title,
		Link: self.Link,
		Tags: []string{}, // ???? I forgor
	}
}

