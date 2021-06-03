package i18n

import (
	"apuscorp.com/core/auth-service/util/constant"
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18n struct {
	Bundle       *i18n.Bundle
	MapLocalizer map[string]*i18n.Localizer
}

func NewI18n() (*I18n, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	mapLocalizer := make(map[string]*i18n.Localizer)
	for lang, path := range constant.LangPathMap {
		bundle.MustLoadMessageFile(path)
		mapLocalizer[lang] = i18n.NewLocalizer(bundle, lang)
	}

	return &I18n{
		Bundle:       bundle,
		MapLocalizer: mapLocalizer,
	}, nil
}

func (r *I18n) MustLocalize(lang string, msgId string, templateData map[string]string) string {
	var localizerPtr *i18n.Localizer
	localizerPtr, ok := r.MapLocalizer[lang]
	if !ok {
		localizerPtr = r.MapLocalizer[constant.DEFAULT_LANG]
	}
	return localizerPtr.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    msgId,
		TemplateData: templateData,
	})
}
