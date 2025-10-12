import EvilIcons from "@expo/vector-icons/EvilIcons";
import { Stack } from "expo-router";
import { useCallback, useState } from "react";
import { ScrollView, TextInput, View } from "react-native";
import { Button } from "@/components/Button";
import { Container } from "@/components/Container";
import { useSearch } from "@/hooks/search/search.hook";
import { mangaMocks } from "@/utils/mocks";
import { MangaPreviewContainer } from "@/components/mangas/preview-container";

export default function SearchPage() {
  const [query, setQuery] = useState<string>();
  const { search, isLoading, results } = useSearch();

  const handleChangeQuery = useCallback((text?: string) => {
    setQuery(text);
  }, []);

  const handleSearch = async () => {
    await search(query);
  };

  return (
    <>
      <Stack.Screen options={{ title: "Search", headerShown: false }} />

      <Container>
        <View className="bg-zinc-900 rounded-md border-2 border-zinc-800 flex flex-row justify-between items-center">
          <TextInput
            value={query}
            onChangeText={handleChangeQuery}
            className=" py-4 px-6  flex-1 placeholder:text-zinc-600 text-zinc-500"
            placeholder="Search..."
          />

          <Button
            className="rounded-l-none !border-0 border-none !bg-zinc-800"
            onPress={handleSearch}
            disabled={isLoading}
          >
            <EvilIcons
              name="search"
              size={24}
              className="stroke-2 mx-auto p-0"
              color="#71717a"
            />
          </Button>
        </View>

        <ScrollView className="flex-1 gap-4 pt-4">
          <View className="flex flex-wrap gap-3 justify-center flex-row w-full flex-1 mt-4">
            {mangaMocks?.map?.((result, idx) => {
              return <MangaPreviewContainer key={idx} manga={result} />;
            })}
          </View>
        </ScrollView>
      </Container>
    </>
  );
}
