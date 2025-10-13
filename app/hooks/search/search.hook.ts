import { useMutation } from "@tanstack/react-query";
import { Toast } from "toastify-react-native";
import { SEARCH_STORE_KEY, useSearchStore } from "@/store/search.store";

export const useSearch = () => {
	const { results, search } = useSearchStore();
	const { mutateAsync, isPending } = useMutation({
		mutationKey: [SEARCH_STORE_KEY, "HOOK"],
		mutationFn: searchMutation,
		onError(error) {
			console.info(error);

			const message = error?.message ?? "Unable to retrieve search results";
			Toast.error(message);
		},
	});

	async function searchMutation(query?: string) {
		if (!query?.length) {
			throw Error(
				"Type something before searching... the scraper can't guess what you want MF",
			);
		}

		const data = await search(query);
		return data;
	}

	return {
		results,
		search: mutateAsync,
		isLoading: isPending,
	};
};
