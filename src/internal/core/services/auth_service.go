package services

import (
	"context"
)

type AuthService interface {
	RequestOTP(ctx context.Context, phone string) error
	VerifyOTP(ctx context.Context, phone, otp string) (token string, err error)
}
