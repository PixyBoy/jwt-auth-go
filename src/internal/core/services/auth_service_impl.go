package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/PixyBoy/jwt-auth-go/internal/core/domain"
	"github.com/PixyBoy/jwt-auth-go/internal/core/ports"
	"github.com/PixyBoy/jwt-auth-go/internal/pkg/util"
	"github.com/rs/zerolog"
)

type AuthServiceImpl struct {
	otpStore    ports.OTPStore
	rateLimiter ports.RateLimiter
	userRepo    ports.UserRepository
	log         zerolog.Logger

	otpDigits          int
	otpTTLSeconds      int
	otpMaxAttempts     int
	otpRateLimitMax    int
	otpRateLimitWindow int
}

func NewAuthService(
	otpStore ports.OTPStore,
	rateLimiter ports.RateLimiter,
	userRepo ports.UserRepository,
	log zerolog.Logger,
	otpDigits int,
	otpTTLSeconds int,
	otpMaxAttempts int,
	otpRateLimitMax int,
	otpRateLimitWindow int,
) AuthService {
	return &AuthServiceImpl{
		otpStore:    otpStore,
		rateLimiter: rateLimiter,
		userRepo:    userRepo,
		log:         log,

		otpDigits:          otpDigits,
		otpTTLSeconds:      otpTTLSeconds,
		otpMaxAttempts:     otpMaxAttempts,
		otpRateLimitMax:    otpRateLimitMax,
		otpRateLimitWindow: otpRateLimitWindow,
	}
}

func (s *AuthServiceImpl) RequestOTP(ctx context.Context, phone string) error {
	allowed, _, err := s.rateLimiter.Allow(phone, s.otpRateLimitMax, s.otpRateLimitWindow)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrRateLimited
	}
	// generate OTP
	otp, err := util.GenerateDigits(s.otpDigits)
	if err != nil {
		return err
	}
	// hash OTP
	h := sha256.Sum256([]byte(otp))
	hash := hex.EncodeToString(h[:])

	// save in redis
	if err := s.otpStore.Save(phone, hash, s.otpTTLSeconds); err != nil {
		return err
	}

	s.log.Info().Str("phone", phone).Str("otp", otp).Msg("OTP generated")
	return nil
}

func (s *AuthServiceImpl) VerifyOTP(ctx context.Context, phone, otp string) (string, error) {
	hash, attempts, exists, err := s.otpStore.Get(phone)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", domain.ErrOTPExpired
	}
	if attempts >= s.otpMaxAttempts {
		_ = s.otpStore.Delete(phone)
		return "", domain.ErrTooManyAttempt
	}

	h := sha256.Sum256([]byte(otp))
	inputHash := hex.EncodeToString(h[:])

	if inputHash != hash {
		if _, err := s.otpStore.IncreaseAttempt(phone); err != nil {
			s.log.Error().Err(err).Msg("increase attempt failed")
		}
		return "", domain.ErrOTPInvalid
	}

	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", err
	}
	if user == nil {
		user, err = s.userRepo.Create(&domain.User{
			Phone: phone,
		})
		if err != nil {
			return "", err
		}
	}

	_ = s.otpStore.Delete(phone)

	token := fmt.Sprintf("verified-user-%d", user.ID)

	return token, nil
}
