type Image = string;

export interface Page {
  id: number;
  content: Image[];
}

export interface Chapter {
  id: number;
  name: string;
  pages: Page[];
}

export interface Manga {
  id?: string;

  image: Image;
  image_url: string
  image_type: string

  url: string
  name: string;
  description?: string;

  tags?: string[] | null
  chapters?: Chapter[];
}
