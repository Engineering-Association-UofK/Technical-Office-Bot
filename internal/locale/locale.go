package locale

import (
	"encoding/json"
	"log/slog"
	"os"
)

type TLocale struct {
	WelcomeMessage string `json:"welcome_message"`
	HelpMessage    string `json:"help_message"`
	FeedbackEmpty  string `json:"feedback_empty"`
	FeedbackThanks string `json:"feedback_thanks"`
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
		slog.Warn("Could not load locale: "+err.Error(), "Language", lang)
		return
	}

	var l TLocale
	if err = json.Unmarshal(data, &l); err != nil {
		slog.Warn("Warning: Could not unmarchal locale: "+err.Error(), "Language", lang)
		return
	}
	lm.locales[lang] = l
}
