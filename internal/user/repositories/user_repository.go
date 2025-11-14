package repositories

import (
	"go-event/internal/user"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uint) (*user.User, error)
	GetAll() ([]*user.User, error)
	Update(user *user.User) error
	Delete(user *user.User) error
	FindByRole(role user.RoleType) ([]*user.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetByID implements UserRepository.
func (r *userRepository) GetByID(id uint) (*user.User, error) {
	var user user.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll implements UserRepository.
func (r *userRepository) GetAll() ([]*user.User, error) {
	var users []*user.User
	if err := r.db.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update implements UserRepository.
func (r *userRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

// Delete implements UserRepository.
func (r *userRepository) Delete(user *user.User) error {
	return r.db.Delete(user).Error
}

// FindByRole implements UserRepository.
func (r *userRepository) FindByRole(role user.RoleType) ([]*user.User, error) {
	var users []*user.User
	if err := r.db.Where("role = ?", role).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
