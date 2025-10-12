import SuccessIcon from "@expo/vector-icons/Feather";
import ErrorIcon from "@expo/vector-icons/Foundation";
import InfoIcon from "@expo/vector-icons/MaterialCommunityIcons";
import { Text, View } from "react-native";

interface Props {
	text1: string;
	text2: string;
	hide?: VoidFunction;
}

export const ErrorToast = ({ text1: title }: Props) => {
	return (
		<View className="bg-zinc-800 border-2 border-zinc-700 rounded-lg flex flex-row items-center gap-3 w-full max-w-[90%] mx-auto p-4 max-h-none break-words">
			<ErrorIcon name="alert" className="my-auto" size={24} color={"#fb7185"} />

			<Text className="text-zinc-200 text-lg font-medium line-clamp-6 break-words w-full text-wrap">
				{title}
			</Text>
		</View>
	);
};

export const InfoToast = ({ text1: title }: Props) => {
	return (
		<View className="bg-zinc-800 border-2 border-zinc-700 rounded-lg flex flex-row items-center gap-3 w-full max-w-[90%] mx-auto p-4 max-h-none break-words">
			<InfoIcon
				name="information-variant-box-outline"
				className="my-auto"
				size={24}
				color={"#a1a1aa"}
			/>

			<Text className="text-zinc-200 text-lg font-medium line-clamp-6 break-words w-full">
				{title}
			</Text>
		</View>
	);
};

export const SuccessToast = ({ text1: title }: Props) => {
	return (
		<View className="bg-zinc-800 border-2 border-zinc-700 rounded-lg flex flex-row items-center gap-3 w-full max-w-[90%] mx-auto p-4 max-h-none">
			<SuccessIcon
				name="check-square"
				className="my-auto"
				size={24}
				color={"#34d399"}
			/>

			<Text className="text-zinc-200 text-lg font-medium line-clamp-6 break-words w-full">
				{title}
			</Text>
		</View>
	);
};
