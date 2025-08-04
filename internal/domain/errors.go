package domain

import "errors"

var (
	ERROR_OUT_OF_STOCK = errors.New("OutOfStock")
	ERROR_NAME_TAKEN   = errors.New("ErrorNameTaken")
	ERROR_ITEM_CODE_TAKEN   = errors.New("ErrorItemCodeTaken")
	ERROR_EMPTY_ARRAY  = errors.New("ErrorEmptyArray")

	ACTION_NOT_ALLOWED   = errors.New("ActionNotAllowed")
	UNEXPECTED_ERROR     = errors.New("UnexpectedError")
	USER_ALREDY_EXIST    = errors.New("UserAlredyExist")
	PARTY_TYPE_NOT_FOUND = errors.New("PartyTypeNotFound")

	NOT_FOUND = errors.New("NoFound")

	ENTITY_NOT_FOUND = errors.New("EntityNoFound")


	PARTY_NOT_FOUND = errors.New("PartyNotFound")

	TYPE_NOT_FOUND = errors.New("TypeNotFound")

	FAIL_TYPE_ASSERTION = errors.New("FailTypeAssertion")

	NIL_POINTER = errors.New("NilPointer")

	BLANK_VALUE = errors.New("BlankValue")

	OVERFLOW_SN                = errors.New("OverflowSn")
	NO_CURRENCY_EXCHANGE_FOUND = errors.New("NoCurrencyExchangeFound")

	DEFAULT_ERRROR = errors.New("Error")

	INVALID_TYPE = errors.New("InvalidType")
)
