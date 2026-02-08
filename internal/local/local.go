package local

import (
	"encoding/json"
	"log"
	"os"
)

type TelegramLocal struct {
	local          string
	WelcomeMessage string `json:"welcome_message"`
	HelpMessage    string `json:"help_message"`
}

func (tl *TelegramLocal) Load() {
	if tl.local == "" {
		tl.local = "en"
	}

	data, err := os.ReadFile("internal/local/" + tl.local + ".json")
	if err != nil {
		log.Fatalln("Failed to load local:", err, " - Will try to use default local.")

		data, err = os.ReadFile("internal/local/en.json")
		if err != nil {
			log.Fatalln("Failed to load default local:", err)
		}
	}

	err = json.Unmarshal(data, &tl)
	if err != nil {
		log.Fatalln("Failed to parse local:", err)
	}
}
