package main

import (
	"flag"

	"github.com/b2r2/Character-Count-Bot/bot"
)

var (
	configPath string
	state      bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.json", "path to config file")
	flag.BoolVar(&state, "set-debug", false, "set debug level\ndefault value: false")
}

func main() {
	flag.Parse()
	bot.Start(state, configPath)
}
