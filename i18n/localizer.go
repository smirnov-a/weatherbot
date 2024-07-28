package i18n

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"weatherbot/internal/logger"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

func Initialize(defaultLanguage string, locales ...string) {
	bundle = i18n.NewBundle(language.Make(defaultLanguage))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, locale := range locales {
		_, err := bundle.LoadMessageFile("i18n/locales/" + locale + ".json")
		if err != nil {
			logger.Logger().Println("Error loading locale file:", err)
		}
	}

	localizer = i18n.NewLocalizer(bundle, defaultLanguage)
}

// SetLocale set current locale
func SetLocale(locale string) {
	localizer = i18n.NewLocalizer(bundle, locale)
}

// Translate the message. if no translation then returns original text
func Translate(messageID string) string {
	translated, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return messageID
	}
	return translated
}
