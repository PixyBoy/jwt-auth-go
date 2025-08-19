package models

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Phone     string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (User) TableName() string { return "users" }
