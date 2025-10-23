import { useQuery } from "@tanstack/react-query";
import { SEARCH_STORE_KEY, useSearchStore } from "@/store/search.store";

export const useManga = (slug?: string) => {
  const { getManga } = useSearchStore();
  const { data, isLoading } = useQuery({
    queryKey: [SEARCH_STORE_KEY, "MANGA", slug],
    queryFn: retrieveManga,
    refetchOnWindowFocus: false,
    enabled: Boolean(slug?.length),
    staleTime: 1000 * 60 * 60,
  });

  async function retrieveManga() {
    if (!slug?.length) {
      throw Error(
        "Slug not found...",
      );
    }

    const data = await getManga(slug);
    return data;
  }

  return {
    isLoading,
    data,
  };
};
