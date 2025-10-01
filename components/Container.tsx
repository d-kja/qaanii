import type { FC, ReactNode } from "react";
import { SafeAreaView } from "react-native";

export const Container: FC<{ children: ReactNode; className?: string }> = ({
  children,
  className = "",
}) => {
  return (
    <SafeAreaView className={`flex flex-1 mx-4 mt-16 ${className}`}>
      {children}
    </SafeAreaView>
  );
};
