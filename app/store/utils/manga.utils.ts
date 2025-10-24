import type { Result } from "../types/search.types";

export const filterManga = (data?: Result | null) => Boolean(data?.slug?.trim?.()?.length)
export const createManga = (data?: Result | null): Result | undefined => {
  if (!data?.slug) {
    return
  }

  return {
    chapters: [],

    ...data,

    // Internal state
    state: {
      error: false,
      loading: false,
      refetching: false,
    },
  }
}
