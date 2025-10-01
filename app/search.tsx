import EvilIcons from "@expo/vector-icons/EvilIcons";
import { Stack } from "expo-router";
import { Text, TextInput, View } from "react-native";
import { Container } from "@/components/Container";

export default function Details() {
  return (
    <>
      <Stack.Screen options={{ title: "Search", headerShown: false }} />

      <Container className="">
        <View className="bg-zinc-900 rounded-md border-2 border-zinc-800 flex flex-row justify-between items-center">
          <TextInput
            className=" py-4 px-6  flex-1 placeholder:text-zinc-600 text-zinc-500"
            placeholder="Search..."
            autoFocus
          />

          <EvilIcons name="search" size={24} color="#52525b" className="px-4" />
        </View>

        <Text className="text-zinc-600">eXAMPLE</Text>
      </Container>
    </>
  );
}
