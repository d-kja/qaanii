import type { Manga } from "./manga.types";

export type Bookmark = Manga

export interface BookmarkStore {
  bookmarks: Bookmark[];

  // Load existing bookmarks into store
  read: () => Promise<void>;

  // I/O
  add: (bookmark: Bookmark) => Promise<void>;
  remove: (bookmark: Bookmark) => Promise<void>;
}
