import { Storage } from "expo-sqlite/kv-store";
import { create } from "zustand";
import { api } from "@/utils/api";
import type { SearchResponseType, SearchStore } from "./types/search.types";

export const SEARCH_STORE_KEY = "@search";
export const useSearchStore = create<SearchStore>((set) => {
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
      console.log("making request");
      const response = await api.get<SearchResponseType>("/search", {
        params: {
          q,
        },
      });
      console.log("finished request");

      const results = response?.data?.mangas ?? [];

      set({
        results,
      });

      await Storage.setItem(SEARCH_STORE_KEY, JSON.stringify(results));
      return results;
    },
  };
});
