package ginadp

import (
	"net/http"

	"github.com/PixyBoy/jwt-auth-go/internal/core/domain"
	"github.com/gin-gonic/gin"
)

func writeDomainError(c *gin.Context, err error) {
	switch err {
	case domain.ErrRateLimited:
		c.JSON(http.StatusTooManyRequests, NewError("OTP_RATE_LIMITED", "try again later"))
	case domain.ErrOTPExpired:
		c.JSON(http.StatusBadRequest, NewError("OTP_EXPIRED", "otp expired or not requested"))
	case domain.ErrOTPInvalid:
		c.JSON(http.StatusBadRequest, NewError("OTP_INVALID", "invalid otp"))
	case domain.ErrTooManyAttempt:
		c.JSON(http.StatusTooManyRequests, NewError("TOO_MANY_ATTEMPTS", "too many attempts"))
	case domain.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, NewError("UNAUTHORIZED", "unauthorized"))
	default:
		c.JSON(http.StatusInternalServerError, NewError("INTERNAL", "internal server error"))
	}
}
