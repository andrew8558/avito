package errors

import "errors"

var (
	ErrAuthFailed          = errors.New("incorrect password")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidHtppMethod   = errors.New("invalid http method")
	ErrNotEnoughMoney      = errors.New("not enough money")
	ErrItemNotFound        = errors.New("item not found")
	ErrInvalidJson         = errors.New("invalid JSON")
	ErrUserDoesNotExist    = errors.New("user requester does not exist")
	ErrReceiverNotFound    = errors.New("receiver not found")
)
