import { create } from "zustand";
import type { DownloadedStore } from "./types/downloaded.types";

export const useDownloaded = create<DownloadedStore>((set, get) => {
  return {
    mangas: [],

    // Load downloaded files
  };
});
