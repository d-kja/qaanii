package utils

import "os"

type Utils struct {}

func (Utils) Envs() map[string]string {
	envs := map[string]string {
		"base_url": os.Getenv("BASE_URL"),
		"port": os.Getenv("PORT"),
	}

	return envs
}
