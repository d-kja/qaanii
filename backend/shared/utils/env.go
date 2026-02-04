package utils

import "os"

func Envs() map[string]string {
	envs := map[string]string{
		"base_url":   os.Getenv("BASE_URL"),
		"broker_url": os.Getenv("BROKER_URL"),
		"port":       os.Getenv("PORT"),
	}

	return envs
}
