import { Storage } from "expo-sqlite/kv-store";
import { create } from "zustand";
import { api } from "@/utils/api";
import type {
  GetMangaResponseType,
  Result,
  SearchResponseType,
  SearchStore,
} from "./types/search.types";
import { createManga, filterManga } from "./utils/manga.utils";

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

    // Utility
    retrieveManga: async (slug: string): Promise<Result | undefined> => {
      const response = await api.get<GetMangaResponseType>(`/manga/${slug}`);
      const manga = response?.data;

      // FIX: scrape the rest to create a new manga instance
      return manga;
    },
    retrieveStoredMangas: async (): Promise<Result[]> => {
      const storedMangas = await Storage.getItem(SEARCH_STORE_KEY);
      if (!storedMangas?.length) {
        return [];
      }

      try {
        const parsedMangas: Result[] = JSON.parse(storedMangas);
        const response = parsedMangas
          ?.map?.(createManga)
          ?.filter?.(filterManga) as Result[];

        return response;
      } catch (error) {
        console.error(
          "[STORE/SEARCH] - Retrieved stored mangas, but couldn't parse.",
          error,
        );
      }

      return [];
    },
    retrieveStoredManga: async (slug: string): Promise<Result | undefined> => {
      const storedManga = await Storage.getItem(`${SEARCH_STORE_KEY}/${slug}`);
      if (!storedManga?.length) {
        return;
      }

      try {
        const parsedManga: Result = JSON.parse(storedManga);
        return parsedManga;
      } catch (error) {
        console.error(
          "[STORE/SEARCH] - Retrieved stored manga, but couldn't parse.",
          error,
        );
      }

      return;
    },
    updateMangas: async (data: Result[]) => {
      set({
        results: data,
      });

      await Storage.setItem(SEARCH_STORE_KEY, JSON.stringify(data));
    },
    updateManga: async (slug: string, data: Result) => {
      await Storage.setItem(
        `${SEARCH_STORE_KEY}/${slug}`,
        JSON.stringify(data),
      );
    },

    // API Actions
    search: async (q?: string) => {
      const { updateMangas, updateManga } = get();
      const response = await api.get<SearchResponseType>("/search", {
        params: {
          q,
        },
      });

      const results = response?.data?.mangas ?? [];

      const procedurePromise: Promise<Result[]> = new Promise(
        (resolve, reject) => {
          console.info("[STORE/SEARCH] - Setting early results");
          set({
            results,
          });

          const steps = async () => {
            try {
              console.info("[STORE/SEARCH] - Updating stored mangas");
              await updateMangas(results);

              for (const result of results) {
                const slug = result?.slug;
                if (!slug?.length) {
                  continue;
                }

                // TODO: Update, this can be slow depending on the search result
                console.info("[STORE/SEARCH] - Updating stored manga:", slug);
                await updateManga(slug, result);
              }

              console.info("[STORE/SEARCH] - Resolving promise");
              resolve(results);
            } catch (err) {
              console.error(
                "[STORE/SEARCH] - Search async procedure failed:",
                err,
              );

              reject(err);
            }
          };

          steps();
        },
      );

      return await procedurePromise;
    },
    getManga: async (slug?: string): Promise<Result | undefined> => {
      if (!slug?.length) {
        return;
      }

      const { retrieveManga, updateManga, retrieveStoredManga } = get();

      const mangas = get().results ?? [];
      const mangaIndex = mangas?.findIndex?.((manga) => manga.slug === slug);

      if (mangaIndex === -1) {
        console.info("[STORE/SEARCH] - In-memory manga not found");
        const manga = await retrieveManga(slug);
        const storedManga = await retrieveStoredManga(slug);

        if (storedManga) {
          console.info("[STORE/SEARCH] - Stored manga found");
          const response: Result = {
            ...storedManga,
            ...manga,
          };

          console.info("[STORE/SEARCH] - Updating stored manga");
          await updateManga(slug, response);
          return response;
        }

        return manga;
      }

      console.info("[STORE/SEARCH] - In-memory manga found");
      let manga = mangas[mangaIndex];
      if (!manga?.slug) {
        console.error(
          "[STORE/SEARCH] - Manga index found, direct indexing not working? This is bad bad...",
        );

        return;
      }

      if (!manga.state) {
        manga = createManga(manga) as Result
      }

      const procedurePromise: Promise<Result | undefined> = new Promise(
        (resolve, reject) => {
          console.info("[STORE/SEARCH] - Pre-updating selected state");
          manga.state.loading = true;

          set({
            selected: manga,
          });

          const steps = async () => {
            try {
              console.info("[STORE/SEARCH] - Searching manga data");
              const newManga = await retrieveManga(slug);

              console.info("[STORE/SEARCH] - Retrieving stored data");
              const storedManga = await retrieveStoredManga(slug);
              console.log(newManga?.chapters)

              const response = createManga({
                ...manga,
                ...storedManga,
                ...newManga,
              }) as Result;

              response.state.loading = false;
              response.state.refetching = false;

              console.info("[STORE/SEARCH] - Closing selected state");
              await updateManga(slug, response);

              set({
                selected: response,
              });

              resolve(manga);
            } catch (err) {
              const error = err as Error;

              console.error(
                "[STORE/SEARCH] - Unable to run procedure, error:",
                error,
              );

              reject(error?.message ?? "Unable to run procedure");
            }
          };

          steps();
        },
      );

      return await procedurePromise;
    },
  };
});
