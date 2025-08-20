package ginadp

import (
	"net/http"
	"strings"

	"github.com/PixyBoy/jwt-auth-go/internal/core/ports"
	"github.com/gin-gonic/gin"
)

const (
	ctxUserIDKey = "auth_user_id"
	ctxPhoneKey  = "auth_phone"
)

func Authz(issuer ports.TokenIssuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, NewError("UNAUTHORIZED", "missing bearer token"))
			c.Abort()
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		uid, phone, err := issuer.Parse(raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, NewError("UNAUTHORIZED", "invalid token"))
			c.Abort()
			return
		}
		c.Set(ctxUserIDKey, uid)
		c.Set(ctxPhoneKey, phone)
		c.Next()
	}
}

func GetAuthUser(c *gin.Context) (int64, string, bool) {
	uidVal, ok1 := c.Get(ctxUserIDKey)
	phVal, ok2 := c.Get(ctxPhoneKey)
	if !ok1 || !ok2 {
		return 0, "", false
	}
	return uidVal.(int64), phVal.(string), true
}
