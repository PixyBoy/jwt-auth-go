package models

import "github.com/PixyBoy/jwt-auth-go/internal/core/domain"

func ToDomainUser(m *User) *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:        int64(m.ID),
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
	}
}

func FromDomainUser(d *domain.User) *User {
	if d == nil {
		return nil
	}
	return &User{
		ID:    uint64(d.ID),
		Phone: d.Phone,
		// CreatedAt: d.CreatedAt,
	}
}
