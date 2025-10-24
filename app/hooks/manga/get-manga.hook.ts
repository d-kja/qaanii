import { useQuery } from "@tanstack/react-query";
import { SEARCH_STORE_KEY, useSearchStore } from "@/store/search.store";

const STAME_TIME = 1000 * 60 * 60 // 60 min

export const useManga = (slug?: string) => {
  const isValidSlug = Boolean(slug?.trim?.()?.length);

  const { getManga, selected } = useSearchStore();
  const { data, isLoading, isRefetching, refetch } = useQuery({
    queryKey: [SEARCH_STORE_KEY, "MANGA", slug],
    queryFn: retrieveManga,
    refetchOnWindowFocus: false,
    enabled: isValidSlug,
    staleTime: STAME_TIME,
  });

  async function retrieveManga() {
    if (!isValidSlug) {
      throw Error("Slug not found...");
    }

    const data = await getManga(slug);
    return data;
  }

  return {
    isLoading: isLoading || isRefetching,
    selected,
    data,

    refresh: refetch,
  };
};
