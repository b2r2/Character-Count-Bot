package bot

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Token    string `json:"Token"`
	Scraping struct {
		Medium string `json:"Medium"`
		Site   string `json:"Site"`
	} `json:"Scraping"`
}

func (c *Config) LoadScrapingConfiguration(configFile string) {
	f, _ := os.Open(configFile)
	decoder := json.NewDecoder(f)
	err := decoder.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}
