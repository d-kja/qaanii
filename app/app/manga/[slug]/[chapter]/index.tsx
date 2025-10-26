import { Stack, useNavigation } from "expo-router";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { useEffect } from "react";
import { ScrollView, Text } from "react-native";
import { Container } from "@/components/Container";
import { useMangaChapter } from "@/hooks/manga-chapter/get-manga-chapter.hook";

type ChapterPageParams = {
	chapter?: string;
};

export default function ChapterPage() {
	const { chapter } = useLocalSearchParams<ChapterPageParams>();
	const navigate = useNavigation();

	const { data: manga, isLoading, refresh } = useMangaChapter(chapter);

	useEffect(() => {
		if (chapter?.length && chapter !== "undefined") {
			return;
		}

		navigate.goBack();
	}, [navigate, chapter]);

	return (
		<>
			<Stack.Screen options={{ title: "Manga Chapter", headerShown: false }} />

			<Container>
				<ScrollView>
					<Text>{chapter}</Text>
				</ScrollView>
			</Container>
		</>
	);
}
