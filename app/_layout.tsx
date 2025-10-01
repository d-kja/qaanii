import { DarkTheme, ThemeProvider } from "@react-navigation/native";
import "../global.css";

import { Stack } from "expo-router";
import { StatusBar } from "react-native";

export default function Layout() {
  return (
    <ThemeProvider value={{
      ...DarkTheme,
      colors: {
        ...DarkTheme.colors,
        background: "#121214",
        text: "#f4f4f5",
      }
    }}>
      <Stack />

      <StatusBar barStyle="light-content" />
    </ThemeProvider>
  );
}
