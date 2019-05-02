package hello

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/tbellembois/gortrocketbot/rocket"
	"golang.org/x/text/language"

	// localization
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
)

type helloPlugin struct {
}

func hello(...string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "hello", PluralCount: 1})
}

func init() {

	// i18n initialization
	bundle = &i18n.Bundle{DefaultLanguage: language.Make(os.Getenv("ROCKETP_HELLO_LANGUAGE"))}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes(LOCALES_EN, "en.toml")
	bundle.MustParseMessageFileBytes(LOCALES_FR, "fr.toml")

	localizer = i18n.NewLocalizer(bundle)

	help := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "help", PluralCount: 1})
	rocket.RegisterPlugin(rocket.Plugin{
		Name:        "hello",            // the rocket command to trigger the action
		CommandFunc: hello,              // command bind function
		Args:        []string{"a", "b"}, // command arguments (not used here, just for example)
		Help:        help,               // command help
		IsAllowed: func(models.User) bool { // everybody can run the plugin
			return true
		},
	})
}
