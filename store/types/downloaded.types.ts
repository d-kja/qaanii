import type { Manga } from "./manga.types";

export type Downloaded = Manga

export interface DownloadedStore {
  mangas: Downloaded[]
}
