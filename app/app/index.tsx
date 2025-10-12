import Entypo from "@expo/vector-icons/Entypo";
import { Link, Stack } from "expo-router";
import {
  Dimensions,
  FlatList,
  ScrollView,
  Text,
  TouchableOpacity,
  View,
} from "react-native";
import { Button } from "@/components/Button";
import { Container } from "@/components/Container";

export default function Home() {
  const { height } = Dimensions.get("window");

  const contentHeight = height - 128 - 48 - 32 - 176 - 16 - 325;

  return (
    <>
      <Stack.Screen options={{ title: "Home", headerShown: false }} />

      <Container className="gap-2">
        <Link href={{ pathname: "/search" }} asChild>
          <Button>
            <View className="flex flex-row justify-between items-center">
              <Text className="text-zinc-600 text-start">Search...</Text>

              <Entypo name="chevron-right" size={18} color="#52525b" />
            </View>
          </Button>
        </Link>

        <ScrollView className="flex-1 gap-4 pt-10">
          <View className="flex-row justify-between">
            <Text className="text-zinc-600 text-lg font-medium">
              You might like...
            </Text>

            <Link href={{ pathname: "/featured" }} asChild>
              <TouchableOpacity className="flex flex-row items-center gap-1 bg-zinc-900 border border-zinc-800 rounded-full px-2.5 py-0.5">
                <Text className="text-zinc-600 text-start">more</Text>

                <Entypo name="chevron-right" size={12} color="#52525b" />
              </TouchableOpacity>
            </Link>
          </View>

          <ScrollView
            horizontal
            contentContainerClassName="gap-4 mt-4"
            className="h-44"
          >
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32"></View>
          </ScrollView>

          <Text className="text-zinc-600 text-lg font-medium pt-4 mt-4">
            Recently updated
          </Text>

          <View className="flex flex-wrap gap-4 justify-center flex-row w-full flex-1 mt-4" style={{ minHeight: contentHeight }}>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>
            <View className="bg-zinc-900 border-2 border-zinc-800 rounded-lg aspect-[9/12] max-w-32 h-44"></View>

          </View>
        </ScrollView>
      </Container>
    </>
  );
}
