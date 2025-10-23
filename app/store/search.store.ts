import { Storage } from "expo-sqlite/kv-store";
import { create } from "zustand";
import { api } from "@/utils/api";
import type {
  GetMangaResponseType,
  Result,
  SearchResponseType,
  SearchStore,
} from "./types/search.types";

export const SEARCH_STORE_KEY = "@search";
export const useSearchStore = create<SearchStore>((set, get) => {
  Storage.getItem(SEARCH_STORE_KEY)
    .then((value) => {
      if (!value?.length) {
        return;
      }

      const results = JSON.parse(value);
      set({
        results,
      });
    })
    .catch((err) => {
      console.error("[ERROR/SEARCH] Unable to load stored results", err);
    });

  return {
    results: [],

    search: async (q?: string) => {
      const response = await api.get<SearchResponseType>("/search", {
        params: {
          q,
        },
      });

      const results = response?.data?.mangas ?? [];

      set({
        results,
      });

      await Storage.setItem(SEARCH_STORE_KEY, JSON.stringify(results));
      return results;
    },

    getManga: async (slug?: string): Promise<Result | undefined> => {
      if (!slug?.length) {
        return;
      }

      const response = await api.get<GetMangaResponseType>(`/manga/${slug}`);
      const manga = response?.data;

      const mangas = get().results;
      const mangaIndex = mangas?.findIndex?.((manga) => manga.slug === slug);
      console.log(mangaIndex)

      if (mangaIndex === -1) {
        return manga;
      }

      const existingManga = structuredClone(mangas[mangaIndex]);
      mangas[mangaIndex] = {
        ...existingManga,
        ...manga,
      };

      set({
        results: mangas,
      });

      await Storage.setItem(SEARCH_STORE_KEY, JSON.stringify(mangas));
      return manga;
    },
  };
});
