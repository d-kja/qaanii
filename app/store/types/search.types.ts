import type { Manga } from "./manga.types";

export type Result = Manga

export type SearchResponseType = {
  mangas: Result[]
}

export interface SearchStore {
  results: Result[];

  search: (query?: string) => Promise<Result[]>;
}
