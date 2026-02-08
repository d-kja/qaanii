package utils

import "os"

func Envs() map[string]string {
	envs := map[string]string{
		"database_url": os.Getenv("DATABASE_URL"),
		"broker_url": os.Getenv("BROKER_URL"),
		"base_url":   os.Getenv("BASE_URL"),
		"redis_url": os.Getenv("REDIS_URL"),
		"port":       os.Getenv("PORT"),
	}

	return envs
}
