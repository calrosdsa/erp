package domain

import "errors"

var (
	ERROR_OUT_OF_STOCK = errors.New("item out of stock")
	ACTION_NOT_ALLOWED = errors.New("ActionNotAllowed")
	UNEXPECTED_ERROR = errors.New("UnexpectedError")
	USER_ALREDY_EXIST  = errors.New("UserAlredyExist")
	PARTY_TYPE_NOT_FOUND = errors.New("PartyTypeNotFound")

	TYPE_NOT_FOUND = errors.New("TypeNotFound")

	FAIL_TYPE_ASSERTION = errors.New("FailTypeAssertion")
)
