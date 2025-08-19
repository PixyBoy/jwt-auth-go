package ginadp

type ReqOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp"   binding:"required"`
}

type VerifyOTPResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewError(code, msg string) ErrorResponse {
	var e ErrorResponse
	e.Error.Code = code
	e.Error.Message = msg
	return e
}
