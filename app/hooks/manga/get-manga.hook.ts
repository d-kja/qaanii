import { useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { Toast } from "toastify-react-native";
import { SEARCH_STORE_KEY, useSearchStore } from "@/store/search.store";
import { queryClient } from "@/utils/query";

const STALE_TIME = 1000 * 60 * 60; // 60 min
const RETRY_TIME = 1000 * 5; // 5s

export const useManga = (slug?: string) => {
  const isValidSlug = Boolean(slug?.trim?.()?.length);

  const { retrieveStoredManga, getManga, updateManga } = useSearchStore();
  const { data, isLoading, isRefetching, refetch } = useQuery({
    queryKey: [SEARCH_STORE_KEY, "MANGA", slug],
    queryFn: retrieveManga,
    refetchOnWindowFocus: false,
    enabled: isValidSlug,
    staleTime: STALE_TIME,
    retryDelay: RETRY_TIME,
    retry: false,
  });

  async function retrieveManga() {
    try {
      if (!slug?.trim?.()?.length || slug === "undefined") {
        throw Error("Slug not found...");
      }

      const stored = await retrieveStoredManga(slug);
      if (stored) {
        console.info("[MANGA] Setting stored manga");
        await queryClient.setQueryData(
          [SEARCH_STORE_KEY, "MANGA", slug],
          stored,
        );
      }

      const data = await getManga(slug);

      console.info("[MANGA] Updating manga");
      await updateManga(slug, data);

      return data;
    } catch (err) {
      const error = err as AxiosError;
      console.info(error);

      const message = "Unable to retrieve manga";
      Toast.error(message);

      if (!slug?.trim?.()?.length || slug === "undefined") {
        return null;
      }

      const stored = await retrieveStoredManga(slug);
      if (stored) {
        console.info(
          "[MANGA] Unable to retrieve updated version, using stored",
        );

        await queryClient.setQueryData(
          [SEARCH_STORE_KEY, "MANGA", slug],
          stored,
        );

        return stored
      }

      return null;
    }
  }

  return {
    isLoading: isLoading || isRefetching,
    data,

    refresh: refetch,
  };
};
