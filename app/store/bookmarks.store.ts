import Storage from "expo-sqlite/kv-store";
import { create } from "zustand";
import type { BookmarkStore } from "./types/bookmark.types";

export const BOOKMARK_STORE_KEY = "@bookmarks";
export const useBookmark = create<BookmarkStore>((set, get) => {
  return {
    bookmarks: [],

    read: async () => {
      const storeBookmarks = await Storage.getItem(BOOKMARK_STORE_KEY);
      if (!storeBookmarks?.length) {
        return;
      }

      try {
        const bookmarks = JSON.parse(storeBookmarks);
        set({ bookmarks });
      } catch (err) {
        console.error("[ERROR]", err);
      }
    },

    add: async (bookmark) => {
      const existingBookmarks = get().bookmarks ?? [];
      const bookmarksMap = new Map(
        existingBookmarks.map((bookmark) => [bookmark.id, bookmark]),
      );

      // Prevent dupes
      bookmarksMap.set(bookmark.id, bookmark);
      const bookmarks = Array.from(bookmarksMap.values());

      await Storage.setItem(BOOKMARK_STORE_KEY, JSON.stringify(bookmarks));
      set({ bookmarks });
    },
    remove: async (bookmark) => {
      const existingBookmarks = get().bookmarks ?? [];
      const bookmarksMap = new Map(
        existingBookmarks.map((bookmark) => [bookmark.id, bookmark]),
      );

      bookmarksMap.delete(bookmark.id);
      const bookmarks = Array.from(bookmarksMap.values());

      await Storage.setItem(BOOKMARK_STORE_KEY, JSON.stringify(bookmarks));
      set({ bookmarks });
    },
  };
});
