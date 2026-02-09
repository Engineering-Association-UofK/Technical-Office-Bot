package locale

import (
	"encoding/json"
	"log"
	"os"
)

type TelegramLocale struct {
	locale         string
	WelcomeMessage string `json:"welcome_message"`
	HelpMessage    string `json:"help_message"`
}

func (tl *TelegramLocale) Load() {
	if tl.locale == "" {
		tl.locale = "en"
	}

	data, err := os.ReadFile("resources/locales/" + tl.locale + ".json")
	if err != nil {
		log.Fatalln("Failed to load locale:", err, " - Will try to use default locale.")

		data, err = os.ReadFile("resources/locales/en.json")
		if err != nil {
			log.Fatalln("Failed to load default locale:", err)
		}
	}

	err = json.Unmarshal(data, &tl)
	if err != nil {
		log.Fatalln("Failed to parse locale:", err)
	}
}
