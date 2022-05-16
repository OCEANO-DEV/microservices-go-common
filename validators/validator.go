package validators

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func NewValidator(
	language string,
) interface{} {
	return func(data interface{}) interface{} {
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
}
