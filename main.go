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
	weaver.Run()
}
