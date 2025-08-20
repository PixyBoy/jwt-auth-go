package ports

import "time"

type TokenIssuer interface {
	Issue(userID int64, phone string, ttl time.Duration) (string, error)
	Parse(token string) (userID int64, phone string, err error)
}
