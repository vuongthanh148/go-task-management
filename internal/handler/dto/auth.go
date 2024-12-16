package dto

import (
	"kn-assignment/internal/core/domain"
	"time"
)

type User struct {
	ID        string      `json:"id"`
	Username  string      `json:"username"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string      `json:"username" example:"employer1"`
	Password string      `json:"password" example:"123456"`
	Role     domain.Role `json:"role" example:"employer"`
}

func (s *CreateUserRequest) ToDomain() domain.CreateUserRequest {
	return domain.CreateUserRequest{
		Username: s.Username,
		Password: s.Password,
		Role:     s.Role,
	}
}

type LoginResponse struct {
	User         User    `json:"user,omitempty"`
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

func (LoginResponse) FromDomain(s domain.LoginResponse) LoginResponse {
	return LoginResponse{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		User: User{
			ID:        s.User.ID,
			Username:  s.User.Username,
			Role:      s.User.Role,
			CreatedAt: s.User.CreatedAt,
			UpdatedAt: s.User.UpdatedAt,
		},
	}
}

type LoginRequest struct {
	Username string `json:"username" example:"employer1"`
	Password string `json:"password" example:"123456"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
