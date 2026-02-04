package entities

import mangav1 "qaanii/mangabuf/gen/manga/v1"

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

func GetStatus(status *string) mangav1.Status {
	switch *status {
	case "completed":
		{
			return mangav1.Status_STATUS_COMPLETED
		}
	case "ongoing":
		{
			return mangav1.Status_STATUS_ONGOING
		}

	default:
		{
			return mangav1.Status_STATUS_UNSPECIFIED
		}
	}
}

func (self Manga) ToProtobuf() mangav1.Manga {
	status := GetStatus(self.Status)
	chapters := []*mangav1.Chapter{}

	if self.Chapters != nil {
		manga_chapters := *self.Chapters

		for _, chapter := range manga_chapters {
			buf_chapter := chapter.ToProtobuf()
			chapters = append(chapters, &buf_chapter)
		}
	}

	return mangav1.Manga{
		Slug:        self.Slug,
		Url:         self.Url,
		Name:        self.Name,
		Description: self.Description,
		Tags:        self.Tags,
		LastUpdate:  nil,
		Image:       self.Image,
		ImageType:   self.ImageType,
		Status:      &status,
		Chapters:    chapters,
	}
}
