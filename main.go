package main

import (
	"github.com/joho/godotenv"
	"github.com/knotseaborg/wikiSearchServer/controller"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	controller.Run()
}
