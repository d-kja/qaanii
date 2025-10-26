import Ionicons from "@expo/vector-icons/Ionicons";
import { Link, Stack, useNavigation } from "expo-router";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { useEffect } from "react";
import { Image, ScrollView, Text, TouchableOpacity, View } from "react-native";
import { Container } from "@/components/Container";
import { useManga } from "@/hooks/manga/get-manga.hook";

type MangaPageParams = {
  slug?: string;
};

export default function MangaPage() {
  const { slug } = useLocalSearchParams<MangaPageParams>();
  const navigate = useNavigation();

  const { data: manga, isLoading, refresh } = useManga(slug);

  useEffect(() => {
    if (slug?.length && slug !== "undefined") {
      return;
    }

    navigate.goBack();
  }, [navigate, slug]);

  const hasImage = !!(manga?.image?.length ?? manga?.image_url?.length);

  const namePlaceholder = isLoading
    ? "Loading name..."
    : "Name not provided...";
  const descriptionPlaceholder = isLoading
    ? "Loading description..."
    : "Description not provided...";

  const uri = manga?.image
    ? `data:image/png;base64,${manga?.image}`
    : manga?.image_url;

  return (
    <>
      <Stack.Screen options={{ title: "Manga", headerShown: false }} />

      <Container>
        <ScrollView>
          <View className="h-[36rem] bg-zinc-900 relative rounded-xl overflow-hidden">
            {hasImage && (
              <Image
                source={{
                  uri,
                }}
                className="h-full w-full object-cover"
              />
            )}

            <TouchableOpacity
              className="absolute top-3 right-2 bg-zinc-900/80 border-2 border-zinc-800 rounded-xl grid place-items-center p-2"
              onPress={() => refresh()}
              disabled={isLoading}
            >
              <Ionicons name="refresh" color="#71717a" size={24} />
            </TouchableOpacity>

            <View className="mt-4 absolute bottom-2.5 bg-zinc-900/90 inset-x-2 rounded-lg px-4 py-2.5 border-2 border-zinc-800">
              <Text className="text-zinc-500 font-semibold text-2xl leading-tight line-clamp-2">
                {manga?.name ?? namePlaceholder}
              </Text>

              <Text className="text-zinc-700 mt-2 line-clamp-6">
                {manga?.description ?? descriptionPlaceholder}
              </Text>
            </View>
          </View>

          <View className="flex-row justify-between items-center mt-6">
            <View className="h-px w-[25%] bg-zinc-800" />

            <View className="flex justify-center items-center">
              <Text className="text-zinc-500">
                {isLoading ? "Thinking..." : "Chapters"}
              </Text>
            </View>

            <View className="h-px w-[25%] bg-zinc-800" />
          </View>

          <View className="flex flex-col gap-2 mt-6">
            {manga?.chapters?.map?.((chapter) => {
              const chapterSlug = chapter.slug;

              return (
                <Link
                  asChild
                  href={`/manga/${chapterSlug}`}
                  key={chapter?.slug}
                >
                  <TouchableOpacity className="flex flex-row justify-between items-center border-2 border-zinc-800 bg-zinc-900 px-4 py-3 rounded-lg max-w-full">
                    <Text className="text-zinc-500 max-w-[70%] line-clamp-2 overflow-hidden">
                      {chapter?.title ?? "-"}
                    </Text>

                    <Text className="text-zinc-700">
                      {chapter?.date ?? "-"}
                    </Text>
                  </TouchableOpacity>
                </Link>
              );
            })}
          </View>
        </ScrollView>
      </Container>
    </>
  );
}
