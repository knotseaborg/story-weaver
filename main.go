package main

import (
	"github.com/joho/godotenv"
	"github.com/knotseaborg/wikiSearchServer/weaver"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	weaver.Run("One of the lead actresses of Avatar 2 appeared to be preparing for another movie.")
}
