package user

import (
	"errors"
	"go-event/internal/notification/email"
	"go-event/pkg/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	ID   uint     `json:"id"`
	Role RoleType `json:"role"`
	jwt.RegisteredClaims
}

type Service interface {
	//for auth
	Register(req RegisterRequest) (*UserResponse, error)
	Login(req LoginRequest) (string, *UserResponse, error)
	GenerateToken(user *User) (string, error)
	//for user
	GetProfile(userID uint) (*UserResponse, error)
	GetAllUsers() ([]UserResponse, error)
	GetUserByID(userID uint) (*UserResponse, error)
	UpdateProfile(userID uint, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(userID uint) error
	GetUsersByRole(role string) ([]UserResponse, error)
	ChangePassword(userID uint, req *ChangePasswordRequest) error
	UpdateRole(userID uint, req *UpdateRoleRequest) (*UserResponse, error)
}

type service struct {
	repo     			Repository
	emailService 	email.Service
	cfg          	*config.Config
}


// GenerateToken implements service.
func (a *service) GenerateToken(user *User) (string, error) {
	// Parse duration from config
	duration, err := time.ParseDuration(a.cfg.JWTExpires)
	if (err != nil) {
		duration = 168 * time.Hour // Default 7 days
	}
	// Create claims with user ID and standard claims
	claims := Claims{
		ID: user.ID,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// Create token with signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token with secret key and return token string
	return token.SignedString([]byte(a.cfg.JWTSecret))
}

// Login implements service.
func (s *service) Login(req LoginRequest) (string,*UserResponse, error) {
	if req.Email == "" || req.Password == "" {
		return "",nil, errors.New("email and password are required")
	}

	users, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "",nil, errors.New("invalid email or password")
		}
		return "",nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(req.Password)); err != nil {
		return "",nil, errors.New("invalid email or password")
	}

	token, err := s.GenerateToken(users)
	if err != nil {
		return "",nil, err
	}

	userResponse := &UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return token,userResponse, nil
}

// Register implements service.
func (s *service) Register(req RegisterRequest) (*UserResponse, error) {
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return nil, errors.New("all fields are required")
	}

	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err	}
	
	newUser := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     RoleParticipant,
	}
	
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}
	
	// Kirim welcome email (async, tidak block jika gagal)
	go func() {
		if err := s.emailService.SendWelcomeEmail(newUser.Email, newUser.Name); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", newUser.Email, err)
		}
	}()
	
	userResponse := &UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  string(newUser.Role),
	}
	return userResponse, nil
}




// GetProfile implements UserService.
func (s *service) GetProfile(userID uint) (*UserResponse, error) {
	users, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	response := &UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// GetAllUsers implements UserService.
func (s *service) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("failed to get users")
	}

	var responses []UserResponse
	for _, u := range users {
		response := UserResponse{
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
func (s *service) GetUserByID(userID uint) (*UserResponse, error) {
	users, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	response := &UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// UpdateProfile implements UserService.
func (s *service) UpdateProfile(userID uint, req *UpdateUserRequest) (*UserResponse, error) {
	users, err := s.repo.GetByID(userID)
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

	if err := s.repo.Update(users); err != nil {
		return nil, errors.New("failed to update user")
	}

	response := &UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return response, nil
}

// DeleteUser implements UserService.
func (s *service) DeleteUser(userID uint) error {
	if err := s.repo.DeleteParticipantsByUserID(userID); err != nil {
        return errors.New("failed to delete related participants: " + err.Error())
    }

	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to get user")
	}

	if err := s.repo.Delete(user); err != nil {
		return errors.New("failed to delete user. " + err.Error())
	}
	return nil
}

// GetUsersByRole implements UserService.
func (s *service) GetUsersByRole(role string) ([]UserResponse, error) {
	users, err := s.repo.FindByRole(RoleType(role))
	if err != nil {
		return nil, errors.New("failed to get users by role")
	}

	var responses []UserResponse
	for _, u := range users {
		response := UserResponse{
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
func (s *service) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to get user")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password" + err.Error())
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)

	if err := s.repo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}


func (s *service) UpdateRole(userID uint, req *UpdateRoleRequest) (*UserResponse, error){
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	if req.Role != "" {
		user.Role = RoleType(req.Role)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, errors.New("failed to update user role: " + err.Error())
	}

	response := &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
	return response, nil
}

func NewService(authRepo Repository, emailService email.Service, cfg *config.Config) Service {
	return &service{
		repo:     		authRepo,
		emailService: emailService,
		cfg:          cfg,
	}
}
