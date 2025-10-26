import { Stack } from "expo-router";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { ScrollView, Text } from "react-native";
import { Container } from "@/components/Container";

type ChapterPageParams = {
	chapter?: string;
};

export default function ChapterPage() {
	const { chapter } = useLocalSearchParams<ChapterPageParams>();

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
