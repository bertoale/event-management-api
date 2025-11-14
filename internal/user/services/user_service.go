package services

import (
	"errors"
	"go-event/internal/user"
	"go-event/internal/user/repositories"
	"go-event/pkg/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(userID uint) (*user.UserResponse, error)
	GetAllUsers() ([]user.UserResponse, error)
	GetUserByID(userID uint) (*user.UserResponse, error)
	UpdateProfile(userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error)
	DeleteUser(userID uint) error
	GetUsersByRole(role string) ([]user.UserResponse, error)
	ChangePassword(userID uint, req *user.ChangePasswordRequest) error
}

type userService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

// GetProfile implements UserService.
func (s *userService) GetProfile(userID uint) (*user.UserResponse, error) {
	users, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	response := &user.UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// GetAllUsers implements UserService.
func (s *userService) GetAllUsers() ([]user.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, errors.New("failed to get users")
	}

	var responses []user.UserResponse
	for _, u := range users {
		response := user.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  string(u.Role),
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// GetUserByID implements UserService.
func (s *userService) GetUserByID(userID uint) (*user.UserResponse, error) {
	users, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	response := &user.UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// UpdateProfile implements UserService.
func (s *userService) UpdateProfile(userID uint, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	users, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	if req.Name != nil {
		users.Name = *req.Name
	}
	if req.Email != nil {
		users.Email = *req.Email
	}

	if err := s.userRepo.Update(users); err != nil {
		return nil, errors.New("failed to update user")
	}

	response := &user.UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// DeleteUser implements UserService.
func (s *userService) DeleteUser(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to get user")
	}

	if err := s.userRepo.Delete(user); err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}

// GetUsersByRole implements UserService.
func (s *userService) GetUsersByRole(role string) ([]user.UserResponse, error) {
	users, err := s.userRepo.FindByRole(user.RoleType(role))
	if err != nil {
		return nil, errors.New("failed to get users by role")
	}

	var responses []user.UserResponse
	for _, u := range users {
		response := user.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  string(u.Role),
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// ChangePassword implements UserService.
func (s *userService) ChangePassword(userID uint, req *user.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to get user")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func NewUserService(userRepo repositories.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}
