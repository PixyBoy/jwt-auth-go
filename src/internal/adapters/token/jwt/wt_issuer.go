package jwtadp

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IssuerHS256 struct {
	secret []byte
}

func NewIssuerHS256(secret string) *IssuerHS256 {
	return &IssuerHS256{secret: []byte(secret)}
}

type claims struct {
	Sub   int64  `json:"sub"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

func (i *IssuerHS256) Issue(userID int64, phone string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := claims{
		Sub:   userID,
		Phone: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(i.secret)
}

func (i *IssuerHS256) Parse(token string) (int64, string, error) {
	t, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return i.secret, nil
	})
	if err != nil {
		return 0, "", err
	}
	if !t.Valid {
		return 0, "", errors.New("invalid token")
	}
	cl, ok := t.Claims.(*claims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}
	return cl.Sub, cl.Phone, nil
}
