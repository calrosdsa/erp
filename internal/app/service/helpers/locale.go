package helpers

import (
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type locale struct {
	bundle          *i18n.Bundle
	defaultLanguage string
	locales         []string
}
type T func(id string) string

type Locale interface {
	MustLocalize(opts ...OptionLocale) (res string)
	GetLang(r *http.Request) string
	Translate(lang string) T
}

func NewLocaleHelper() Locale {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// locales := viper.GetStringSlice("locales")
	locales := []string{"en", "es"}
	for _, locale := range locales {
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/util.%s.toml", locale))
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/util.%s.toml", locale))
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/active.%s.toml", locale))
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/active.%s.toml", locale))
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/template.%s.toml", locale))
		bundle.MustLoadMessageFile(fmt.Sprintf("../../asset/locale/template.%s.toml", locale))

	}
	// lo := locales[0]
	return &locale{
		bundle:          bundle,
		defaultLanguage: locales[0],
		locales:         locales,
	}
}
func (l *locale) Translate(lang string) T {
	return func(id string) string {
		localizer := i18n.NewLocalizer(l.bundle, lang)
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: id,
			},
		})
	}
}

func (l *locale) GetLang(r *http.Request) string {
	lang := r.Header["Accept-Language"][0]
	for _, locale := range l.locales {
		if locale == lang {
			return lang
		}
	}
	return l.defaultLanguage
}

func (l *locale) MustLocalize(opts ...OptionLocale) (res string) {
	options := OptionsLocale.Apply(opts...)
	if options.Lang == "" {
		options.Lang = l.defaultLanguage
	}
	localizer := i18n.NewLocalizer(l.bundle, options.Lang)
	res = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    options.ID,
			One:   options.One,
			Other: options.Other,
		},
		TemplateData: options.Template,
	})
	return
}

type optionsLocale struct {
	Lang     string
	ID       string
	One      string
	Other    string
	Template interface{}
}

type OptionLocale func(o *optionsLocale)

var OptionsLocale optionsLocale

func (o *optionsLocale) WithLang(lang string) OptionLocale {
	return func(o *optionsLocale) {
		o.Lang = lang
	}
}

func (o *optionsLocale) WithID(id string) OptionLocale {
	return func(o *optionsLocale) {
		o.ID = id
	}
}

func (o *optionsLocale) WithOne(one string) OptionLocale {
	return func(o *optionsLocale) {
		o.One = one
	}
}

func (o *optionsLocale) WithOther(other string) OptionLocale {
	return func(o *optionsLocale) {
		o.Other = other
	}
}

func (o *optionsLocale) WithTemplate(template interface{}) OptionLocale {
	return func(o *optionsLocale) {
		o.Template = template
	}
}

func (o *optionsLocale) Apply(opts ...OptionLocale) optionsLocale {
	options := optionsLocale{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
