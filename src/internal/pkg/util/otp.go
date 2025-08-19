package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateDigits(n int) (string, error) {
	max := big.NewInt(10)
	otp := ""
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num.Int64())
	}
	return otp, nil
}
