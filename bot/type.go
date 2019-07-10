package bot

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Token  string `json:"Token_Test"`
	Medium string `json:"Medium"`
	Site   struct {
		Login    string `json:"Login"`
		Password string `json:"Password"`
		Domain   string `json:Domain"`
		URL      string `json:"URL"`
	} `json:"Site"`
}

func (c *Config) LoadScrapingConfiguration(configFile string) {
	f, _ := os.Open(configFile)
	decoder := json.NewDecoder(f)
	err := decoder.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}

type WPResponse struct {
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
}
