package validators

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type validators struct {
	trans ut.Translator
}

var language string

func NewValidator(
	language string,
) *validators {
	var uni *ut.UniversalTranslator
	var trans ut.Translator

	switch language {
	case "en":
		uni = ut.New(en.New())
	case "pt_BR":
		uni = ut.New(pt_BR.New())
	default:
		uni = ut.New(en.New())
	}

	trans, _ = uni.GetTranslator(language)

	return &validators{
		trans: trans,
	}
}

func (v *validators) Validate(data interface{}) interface{} {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator(language)
	validate := validator.New()

	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(trans))
		}

		return errors
	}

	return nil
}
