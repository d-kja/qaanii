import type { Manga } from "./manga.types";

export type Result = Manga;

export type SearchResponseType = {
  mangas: Result[];
};

export type GetMangaResponseType = Result;

export interface SearchStore {
  selected?: Result;
  results: Result[];

  search: (query?: string) => Promise<Result[]>;
  getManga: (slug?: string) => Promise<Result | undefined>;

  retrieveStoredManga: (slug: string) => Promise<Result | undefined>;
  retrieveManga: (slug: string) => Promise<Result | undefined>;
  updateManga: (slug: string, data: Result) => Promise<void>;

  updateMangas: (data: Result[]) => Promise<void>;
}
