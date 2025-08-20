package mysqladp

import (
	"github.com/PixyBoy/jwt-auth-go/internal/adapters/db/mysql/models"
	"github.com/PixyBoy/jwt-auth-go/internal/core/domain"
	"github.com/PixyBoy/jwt-auth-go/internal/core/ports"
	"gorm.io/gorm"
)

type UserRepoGorm struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepoGorm{db: db}
}

func (r *UserRepoGorm) FindByID(id int64) (*domain.User, error) {
	var m models.User
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	u := models.ToDomainUser(&m)
	if u == nil {
		println("ToDomainUser returned nil for user")
	}
	return u, nil
}

func (r *UserRepoGorm) FindByPhone(phone string) (*domain.User, error) {
	var m models.User
	if err := r.db.Where("phone = ?", phone).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	u := models.ToDomainUser(&m)
	if u == nil {
		println("ToDomainUser returned nil for user with phone:", phone)
	}
	return u, nil
}

func (r *UserRepoGorm) Create(u *domain.User) (*domain.User, error) {
	m := models.FromDomainUser(u)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return models.ToDomainUser(m), nil
}

func (r *UserRepoGorm) List(search string, page, perPage int) ([]domain.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 200 {
		perPage = 20
	}
	var ms []models.User
	q := r.db.Model(&models.User{})
	if search != "" {
		q = q.Where("phone LIKE ?", "%"+search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Order("id DESC").
		Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	res := make([]domain.User, 0, len(ms))
	for _, m := range ms {
		res = append(res, *models.ToDomainUser(&m))
	}
	return res, total, nil
}
