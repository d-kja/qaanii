type Image = any; // ???

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
  id: string;

  cover: Image;
  name: string;
  description?: string;

  chapters: Chapter[];
}
