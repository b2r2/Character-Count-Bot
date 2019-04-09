package main

import (
	"Character-Count-Bot/bot"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Token string
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	checkError(err)

	b, e := bot.Init(configuration.Token, true)
	checkError(e)

	bot.Start(b)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
