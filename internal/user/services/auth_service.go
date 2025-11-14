package services

import (
	"errors"
	"go-event/internal/notification/email"
	"go-event/internal/user"
	"go-event/internal/user/repositories"
	"go-event/pkg/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	ID   uint          `json:"id"`
	Role user.RoleType `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(req user.RegisterRequest) (*user.UserResponse, error)
	Login(req user.LoginRequest) (string, *user.UserResponse, error)
	GenerateToken(user *user.User) (string, error)
}

type authService struct {
	authRepo     repositories.AuthRepository
	emailService email.Service
	cfg          *config.Config
}

// GenerateToken implements AuthService.
func (a *authService) GenerateToken(user *user.User) (string, error) {
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

// Login implements AuthService.
func (a *authService) Login(req user.LoginRequest) (string,*user.UserResponse, error) {
	if req.Email == "" || req.Password == "" {
		return "",nil, errors.New("email and password are required")
	}

	users, err := a.authRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "",nil, errors.New("invalid email or password")
		}
		return "",nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(req.Password)); err != nil {
		return "",nil, errors.New("invalid email or password")
	}

	token, err := a.GenerateToken(users)
	if err != nil {
		return "",nil, err
	}

	userResponse := &user.UserResponse{
		ID:    users.ID,
		Name:  users.Name,
		Email: users.Email,
		Role:  string(users.Role),
	}
	return token,userResponse, nil
}

// Register implements AuthService.
func (a *authService) Register(req user.RegisterRequest) (*user.UserResponse, error) {
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return nil, errors.New("all fields are required")
	}

	existingUser, _ := a.authRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err	}
	
	newUser := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     user.RoleParticipant,
	}
	
	if err := a.authRepo.Register(newUser); err != nil {
		return nil, err
	}
	
	// Kirim welcome email (async, tidak block jika gagal)
	go func() {
		if err := a.emailService.SendWelcomeEmail(newUser.Email, newUser.Name); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", newUser.Email, err)
		}
	}()
	
	userResponse := &user.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  string(newUser.Role),
	}
	return userResponse, nil
}

// GenerateToken implements AuthService.

func NewAuthService(authRepo repositories.AuthRepository, emailService email.Service, cfg *config.Config) AuthService {
	return &authService{
		authRepo:     authRepo,
		emailService: emailService,
		cfg:          cfg,
	}
}
