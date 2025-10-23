import { Link } from "expo-router";
import type { FC } from "react";
import { Image, Text, TouchableOpacity, View } from "react-native";
import type { Manga } from "@/store/types/manga.types";

interface Props {
  manga?: Manga;
}

export const MangaPreviewContainer: FC<Props> = ({ manga }) => {
  return (
    <Link href={`/manga/${manga?.slug}`} asChild>
      <TouchableOpacity className="max-w-[8.5rem] items-center justify-center px-2 py-1 flex gap-2 bg-zinc-900 border-2 border-zinc-800 rounded-lg">
        <View className="aspect-[9/12] max-w-[7.75rem] h-44 overflow-hidden bg-zinc-800/40 rounded-md">
          <Image
            source={{
              uri: manga?.image
                ? `data:image/png;base64,${manga?.image}`
                : manga?.image_url,
            }}
            className="flex-1  rounded-md"
          />
        </View>

        <View className="flex flex-col gap-2 pb-2">
          <Text className="line-clamp-1 text-zinc-500 text-xs">
            {manga?.name ?? "Not found"}
          </Text>
        </View>
      </TouchableOpacity>
    </Link>
  );
};
