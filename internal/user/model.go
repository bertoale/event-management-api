package user

import "time"

type RoleType string

const (
	RoleAdmin       RoleType = "admin"
	RoleOrganizer   RoleType = "organizer"
	RoleParticipant RoleType = "participant"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:191"`
	Password  string    `json:"-"` // jangan dikirim ke response
	Role      RoleType  `json:"role" gorm:"type:enum('admin','organizer','participant')"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  string(u.Role),
	}
}

// 📩 Request structs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin organizer participant"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin organizer participant"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" form:"old_password" validate:"required"`
    NewPassword string `json:"new_password" form:"new_password" validate:"required,min=6"`
}


// 📤 Response structs
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Participant struct {
	ID     uint
	UserID uint
}