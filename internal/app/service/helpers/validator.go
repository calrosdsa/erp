package helpers

 import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
 )


 type ValidatorHelper struct {
	Validate *validator.Validate
	Uni *ut.UniversalTranslator
 }

func NewValidator() *ValidatorHelper{
	en := en.New()
	es := es.New()
	uni := ut.New(en, en,es)
	transEn, _ := uni.GetTranslator("en")
	transEs, _ := uni.GetTranslator("es")
	validate := validator.New()
	registerTranslationEn("required",transEn,validate)
	registerTranslationEs("required",transEs,validate)
	return &ValidatorHelper{
		Validate: validate,
		Uni: uni,
	}
}

func registerTranslationEn(key string,trans ut.Translator,validate *validator.Validate){
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "This field is required.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required")
		return t
	})

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "Please enter a valid email address.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email")
		return t
	})

	validate.RegisterTranslation("e164", trans, func(ut ut.Translator) error {
		return ut.Add("e164", "Please enter a valid phone number.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("e164")
		return t
	})
	
}

func registerTranslationEs(key string,trans ut.Translator,validate *validator.Validate){
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "Este campo es obligatorio.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "Por favor, ingresa una dirección de correo electrónico válida.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email")
		return t
	})

	validate.RegisterTranslation("e164", trans, func(ut ut.Translator) error {
		return ut.Add("e164", "Por favor, ingresa un número de teléfono válido.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("e164")
		return t
	})
}