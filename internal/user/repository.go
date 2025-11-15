package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	//for auth
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	//for profile
	GetByID(id uint) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(user *User) error
	FindByRole(role RoleType) ([]*User, error)
	DeleteParticipantsByUserID(userID uint) error
}

// Tambahkan model Participant untuk query delete
// Model minimal agar bisa digunakan untuk delete



type repository struct {
	db *gorm.DB
}

// FindByEmail implements Repository.
func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create implements Repository.
func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// GetByID implements UserRepository.
func (r *repository) GetByID(id uint) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll implements UserRepository.
func (r *repository) GetAll() ([]*User, error) {
	var users []*User
	if err := r.db.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update implements UserRepository.
func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

// Delete implements UserRepository.
func (r *repository) Delete(user *User) error {
	return r.db.Delete(user).Error
}

// FindByRole implements UserRepository.
func (r *repository) FindByRole(role RoleType) ([]*User, error) {
	var users []*User
	if err := r.db.Where("role = ?", role).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteParticipantsByUserID menghapus semua participant dengan user_id tertentu
func (r *repository) DeleteParticipantsByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&Participant{}).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
