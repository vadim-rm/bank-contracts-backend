package domain

import "errors"

var (
	ErrNotFound            = errors.New("ERR_NOT_FOUND")
	ErrAccountNumberEmpty  = errors.New("ERR_ACCOUNT_NUMBER_EMPTY")
	ErrInvalidTargetStatus = errors.New("ERR_INVALID_TARGET_STATUS")
	ErrActionNotPermitted  = errors.New("ERR_ACTION_NOT_PERMITTED")
	ErrWrongAccountStatus  = errors.New("ERR_WRONG_ACCOUNT_STATUS")
)
