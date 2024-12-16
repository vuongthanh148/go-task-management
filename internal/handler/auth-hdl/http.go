package authhdl

import (
	"kn-assignment/internal/core/port"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type handler struct {
	svc port.AuthService
}

func New(authService port.AuthService) Handler {
	return &handler{
		svc: authService,
	}
}
