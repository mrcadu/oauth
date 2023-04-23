package locale

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Setup() {
	bundle = i18n.NewBundle(language.BrazilianPortuguese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("locale/en-US.json")
	bundle.MustLoadMessageFile("locale/pt-BR.json")
}

func GetMessage(messageId string, lang string, accept string) string {
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	localizeConfig := i18n.LocalizeConfig{
		MessageID: messageId,
	}
	localize, err := localizer.Localize(&localizeConfig)
	if err != nil {
		panic(err)
	}
	return localize
}

func GetMessageLocaleFromRequest(messageId string, ctx *gin.Context, params map[string]string) string {
	lang := ctx.PostForm("lang")
	accept := ctx.Request.Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	localizeConfig := i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: params,
	}
	localize, err := localizer.Localize(&localizeConfig)
	if err != nil {
		panic(err)
	}
	return localize
}
