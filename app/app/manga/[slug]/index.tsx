import { Stack } from "expo-router";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { Image, Text } from "react-native";
import { Container } from "@/components/Container";
import { useManga } from "@/hooks/manga/get-manga.hook";

type MangaPageParams = {
  slug?: string
}

export default function MangaPage() {
  const { slug } = useLocalSearchParams<MangaPageParams>()
  const { data, isLoading } = useManga(slug)
  console.log(data, isLoading)

  return (
    <>
      <Stack.Screen options={{ title: "Manga", headerShown: false }} />

      <Container>
        <Image
          source={{
              uri: data?.image
                ? `data:image/png;base64,${data?.image}`
                : data?.image_url,
          }}
        />
        <Text className="text-zinc-100">{data?.name}</Text>

        {data?.chapters?.map?.(chapter => 
          <Text key={chapter?.slug} className="text-zinc-100">{chapter?.title ?? "-"}</Text>
        )}
      </Container>
    </>
  );
}
