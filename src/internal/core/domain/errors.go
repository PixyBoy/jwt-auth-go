package domain

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidInput   = errors.New("invalid input")
	ErrRateLimited    = errors.New("rate limited")
	ErrOTPExpired     = errors.New("otp expired")
	ErrOTPInvalid     = errors.New("otp invalid")
	ErrTooManyAttempt = errors.New("too many attempts")
)
