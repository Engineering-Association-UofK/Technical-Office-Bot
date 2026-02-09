package locale

import (
	"encoding/json"
	"log"
	"os"
)

type TLocale struct {
	WelcomeMessage string `json:"welcome_message"`
	HelpMessage    string `json:"help_message"`
}

type LocaleManager struct {
	locales map[string]TLocale
}

func NewLocaleManager() *LocaleManager {
	lm := &LocaleManager{
		locales: make(map[string]TLocale),
	}

	lm.load("en")
	lm.load("ar")
	return lm
}

func (lm *LocaleManager) Get(lang string) TLocale {
	if val, ok := lm.locales[lang]; ok {
		return val
	}
	return lm.locales["en"]
}

func (lm *LocaleManager) load(lang string) {
	data, err := os.ReadFile("resources/locales/" + lang + ".json")
	if err != nil {
		log.Printf("Warning: Could not load locale %s: %v", lang, err)
		return
	}

	var l TLocale
	if err = json.Unmarshal(data, &l); err != nil {
		log.Printf("Warning: Could not unmarchal locale %s: %v", lang, err)
		return
	}
	lm.locales[lang] = l
}
