import { DarkTheme, ThemeProvider } from "@react-navigation/native";
import "../global.css";

import { QueryClientProvider } from "@tanstack/react-query";
import { Stack } from "expo-router";
import { StatusBar } from "react-native";
import ToastManager from "toastify-react-native";
import { queryClient } from "@/utils/query";
import { toastConfig } from "@/utils/toast.config";

export default function Layout() {
  return (
    <ThemeProvider
      value={{
        ...DarkTheme,
        colors: {
          ...DarkTheme.colors,
          background: "#121214",
          text: "#f4f4f5",
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <Stack />

        <StatusBar barStyle="light-content" />
        <ToastManager config={toastConfig} animationStyle="slide" position="bottom" />
      </QueryClientProvider>
    </ThemeProvider>
  );
}
