import { useQuery } from "@tanstack/react-query";
import { SEARCH_STORE_KEY, useSearchStore } from "@/store/search.store";
import { queryClient } from "@/utils/query";

const STAME_TIME = 1000 * 60 * 60; // 60 min

export const useManga = (slug?: string) => {
  const isValidSlug = Boolean(slug?.trim?.()?.length);

  const { retrieveStoredManga, getManga, updateManga } = useSearchStore();
  const { data, isLoading, isRefetching, refetch } = useQuery({
    queryKey: [SEARCH_STORE_KEY, "MANGA", slug],
    queryFn: retrieveManga,
    refetchOnWindowFocus: false,
    enabled: isValidSlug,
    staleTime: STAME_TIME,
  });

  async function retrieveManga() {
    if (!slug?.trim?.()?.length) {
      throw Error("Slug not found...");
    }

    console.info("[MANGA] Retrieving stored manga");
    const stored = await retrieveStoredManga(slug);
    if (stored) {
      console.info("[MANGA] Setting stored manga");
      await queryClient.setQueryData([SEARCH_STORE_KEY, "MANGA", slug], stored);
    }

    console.info("[MANGA] Updating manga");
    const data = await getManga(slug);

    console.info("[MANGA] Persisting manga");
    await updateManga(slug, data);

    return data;
  }

  return {
    isLoading: isLoading || isRefetching,
    data,

    refresh: refetch,
  };
};
