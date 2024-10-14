package domain

import "errors"

var (
	ErrNotFound            = errors.New("ERR_NOT_FOUND")
	ErrAccountNumberEmpty  = errors.New("ERR_ACCOUNT_NUMBER_EMPTY")
	ErrInvalidTargetStatus = errors.New("ERR_INVALID_TARGET_STATUS")
	ErrActionNotPermitted  = errors.New("ERR_ACTION_NOT_PERMITTED")
	ErrWrongAccountStatus  = errors.New("ERR_WRONG_ACCOUNT_STATUS")
	ErrInvalidCredentials  = errors.New("ERR_INVALID_CREDENTIALS")
	ErrUnknown             = errors.New("ERR_UNKNOWN")
	ErrUnauthenticated     = errors.New("ERR_UNAUTHENTICATED")
	ErrBadRequest          = errors.New("ERR_BAD_REQUEST")
)
