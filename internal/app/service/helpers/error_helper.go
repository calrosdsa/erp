package helpers

import (
	"encoding/json"
	"erp/api/common"
	"erp/internal/domain"
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

type ErrorHelper interface {
	HumaCustomError(lng string, err error, args ...interface{}) huma.StatusError
	ErrorMessage(opts ...OptionLocale) error
}

type errorHelper struct {
	locale Locale
}

func NewErrorHelper(locale Locale) ErrorHelper {
	return &errorHelper{
		locale: locale,
	}
}

func (h *errorHelper)ErrorMessage(opts ...OptionLocale) error{
	message := h.locale.MustLocalize(opts...)
	errorMessage:= common.ErrorMessage{
		Message: message,
	}
	jsonData, err := json.Marshal(&errorMessage)
	if err != nil {
		return domain.DEFAULT_ERRROR
	}

	// Convert the byte slice to a string and print it
	jsonString := string(jsonData)
	return errors.New(jsonString)
}

func (h *errorHelper) HumaCustomError(lng string, err error, args ...interface{}) huma.StatusError {
	var customErrorMsg string
    var message common.ErrorMessage

    // Attempt to unmarshal the error into a structured message
    if jsonErr := json.Unmarshal([]byte(err.Error()), &message); jsonErr == nil {
        return huma.Error400BadRequest(message.Message)
    }

    // If there are additional arguments, check for a custom error message
    if len(args) > 0 {
        if value, ok := args[0].(string); ok {
            customErrorMsg = value
        }
    }

    // Translate error messages based on language
    t := h.locale.Translate(lng)

    // Define a map of domain errors to translated messages
    errorTranslations := map[error]string{
        domain.ERROR_OUT_OF_STOCK:        t("Error.OutOfStock"),
        domain.ERROR_ITEM_CODE_TAKEN:     t("Error.ErrorItemCodeTaken"),
        domain.OVERFLOW_SN:               t("Error.OverflowSn"),
        domain.ACTION_NOT_ALLOWED:        t("Error.ActionNotAllowed"),
    }

    // Return the corresponding error message based on the domain error
    if translatedMessage, exists := errorTranslations[err]; exists {
        return huma.Error400BadRequest(translatedMessage)
    }

    // If a custom error message is provided, return it
    if customErrorMsg != "" {
        return huma.Error400BadRequest(t(customErrorMsg))
    }

    // Default error message if no specific case is matched
    return huma.Error400BadRequest(t("Error.FailToFetchData"))
}
