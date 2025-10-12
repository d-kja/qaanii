package main

import (
	"server/internals/http"
	"server/internals/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.LOGGER.ERROR.Fatal(err)
	}

	handler := http.HttpHandler{}
	handler.Run()
}

func init() {
	utils.Logger{}.Init()
}
