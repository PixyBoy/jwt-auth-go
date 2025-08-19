package ports

import "github.com/PixyBoy/jwt-auth-go/internal/core/domain"

type UserRepository interface {
	FindByID(id int64) (*domain.User, error)
	FindByPhone(phone string) (*domain.User, error)
	Create(u *domain.User) (*domain.User, error)
	List(search string, page, perPage int) ([]domain.User, int64, error)
}
