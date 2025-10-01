import { forwardRef } from "react";
import {
  TouchableOpacity,
  type TouchableOpacityProps,
  type View,
} from "react-native";

type ButtonProps = {} & TouchableOpacityProps;

export const Button = forwardRef<View, ButtonProps>(
  ({ children, className = "", ...touchableProps }, ref) => {
    return (
      <TouchableOpacity
        ref={ref}
        {...touchableProps}
        className={`bg-zinc-900 py-4 px-6 rounded-md border-2 border-zinc-800 ${className}`}
      >
        {children}
      </TouchableOpacity>
    );
  },
);

Button.displayName = "Button";
