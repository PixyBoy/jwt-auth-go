package ginadp

import (
	"net/http"

	"github.com/PixyBoy/jwt-auth-go/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// @Summary     Request OTP
// @Description Generates OTP and stores hash in Redis (rate limited).
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       payload body ReqOTPRequest true "Phone"
// @Success     204
// @Failure     400 {object} ErrorResponse
// @Failure     429 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /v1/auth/otp/request [post]
func RequestOTPHandler(auth services.AuthService, log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ReqOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewError("INVALID_INPUT", "phone is required"))
			return
		}
		if err := auth.RequestOTP(c.Request.Context(), req.Phone); err != nil {
			writeDomainError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// @Summary     Verify OTP
// @Description Verifies OTP; creates user if not exists; returns token placeholder.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       payload body VerifyOTPRequest true "Phone & OTP"
// @Success     200 {object} VerifyOTPResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Failure     429 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /v1/auth/otp/verify [post]
func VerifyOTPHandler(auth services.AuthService, log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewError("INVALID_INPUT", "phone and otp are required"))
			return
		}
		token, err := auth.VerifyOTP(c.Request.Context(), req.Phone, req.OTP)
		if err != nil {
			writeDomainError(c, err)
			return
		}
		c.JSON(http.StatusOK, VerifyOTPResponse{
			Token: token,
			User:  map[string]any{"phone": req.Phone},
		})
	}
}
