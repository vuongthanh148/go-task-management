package authsvc

import (
	"kn-assignment/internal/core/port"
)

type service struct {
	repo port.AuthRepository
}

func New(authRepository port.AuthRepository) port.AuthService {
	return &service{repo: authRepository}
}
