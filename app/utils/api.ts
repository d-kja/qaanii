import axios from "axios";

const baseURL = process.env.EXPO_PUBLIC_API_URL;
if (!baseURL) {
  throw new Error("Missing API_URL environment variable");
}

export const api = axios.create({
  baseURL,
});
