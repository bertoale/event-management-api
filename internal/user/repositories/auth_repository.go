package repositories

import (
	"go-event/internal/user"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*user.User, error)
	Register(user *user.User) error
}

type authRepository struct {
	db *gorm.DB
}

// FindByEmaile implements Repository.
func (r *authRepository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil 
}

// Register implements Repository.
func (r *authRepository) Register(user *user.User) error {
	return r.db.Create(user).Error
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
