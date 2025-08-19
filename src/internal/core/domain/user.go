package domain

import "time"

type User struct {
	ID        int64
	Phone     string
	CreatedAt time.Time
}
