import { Stack } from "expo-router";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { Image, Text, TouchableOpacity, View } from "react-native";
import { Container } from "@/components/Container";
import { useManga } from "@/hooks/manga/get-manga.hook";
import Ionicons from "@expo/vector-icons/Ionicons";

type MangaPageParams = {
  slug?: string;
};

export default function MangaPage() {
  const { slug } = useLocalSearchParams<MangaPageParams>();
  const { selected, data, isLoading, refresh } = useManga(slug);

  const loading = isLoading || selected?.state?.loading;
  console.log(selected?.chapters, data?.chapters);
  return (
    <>
      <Stack.Screen options={{ title: "Manga", headerShown: false }} />

      <Container>
        <View className="h-[36rem] bg-zinc-900 relative rounded-xl overflow-hidden">
          {!loading && (
            <Image
              source={{
                uri: data?.image
                  ? `data:image/png;base64,${data?.image}`
                  : data?.image_url,
              }}
              className="h-full w-full object-cover opacity-60" // Can't blind myself while working...
            />
          )}

          <TouchableOpacity
            className="absolute top-3 right-2 bg-zinc-900/80 border-2 border-zinc-800 rounded-xl grid place-items-center p-2"
            onPress={() => refresh()}
            disabled={loading}
          >
            <Ionicons name="refresh" color="#52525b" size={24} />
          </TouchableOpacity>

          <View className="mt-4 absolute bottom-2.5 bg-zinc-900/90 inset-x-2 rounded-lg px-4 py-2.5 border-2 border-zinc-800">
            <Text className="text-zinc-500 font-semibold text-2xl leading-tight line-clamp-2">
              {data?.name}
            </Text>

            <Text className="text-zinc-700 mt-2 line-clamp-6">
              {data?.description ?? "Description not provided..."}
            </Text>
          </View>
        </View>
  
        <View className="h-px w-full bg-zinc-800 mt-6" />

        <View className="flex flex-col gap-2 mt-6">
          {selected?.chapters?.map?.((chapter) => (
            <View key={chapter?.slug} className="flex flex-row justify-between items-center border-2 border-zinc-800 bg-zinc-900 px-4 py-3 rounded-lg">
              <Text className="text-zinc-500">{chapter?.title ?? "-"}</Text>

              <Text className="text-zinc-700">{chapter?.date ?? "-"}</Text>
            </View>
          ))}
        </View>
      </Container>
    </>
  );
}
