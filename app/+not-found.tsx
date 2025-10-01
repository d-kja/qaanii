import Entypo from "@expo/vector-icons/Entypo";
import { Stack, useNavigation } from "expo-router";
import { Text, TouchableOpacity } from "react-native";
import { Container } from "@/components/Container";

export default function NotFoundScreen() {
  const navigate = useNavigation()

  const handleReturn = () => {
    navigate.goBack()
  }
  return (
    <>
      <Stack.Screen options={{ title: "Oops!", headerShown: false }} />

      <Container className="justify-center items-center -mt-12 gap-2">
        <Text className={"text-4xl font-bold text-zinc-400"}>404</Text>
        <Text className={"text-xl font-bold text-zinc-400"}>
          Page not found!
        </Text>

        <TouchableOpacity onPress={handleReturn} className="mt-6 py-2 px-6 flex flex-row items-center justify-center gap-2">
          <Entypo name="chevron-left" size={14} color="#52525b" className="" />
          <Text className={"text-base text-zinc-500"}>Go back!</Text>
        </TouchableOpacity>
      </Container>
    </>
  );
}
