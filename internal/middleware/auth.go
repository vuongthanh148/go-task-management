package middleware

import (
	"kn-assignment/internal/constant"
	"kn-assignment/internal/core/domain"
	errors "kn-assignment/internal/core/error"
	"kn-assignment/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.HandleError(c, constant.ErrCodeUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := util.ValidateToken(tokenStr)
		if err != nil {
			errors.HandleError(c, constant.ErrCodeUnauthorized)
			return
		}

		c.Set("username", claims.Username)
		c.Set("userId", claims.Id)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(requiredRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			errors.HandleError(c, constant.ErrCodeForbidden)
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			errors.HandleError(c, constant.ErrCodeForbidden)
			return
		}
		if domain.Role(roleStr) != requiredRole {
			errors.HandleError(c, constant.ErrCodeForbidden)
			return
		}

		c.Next()
	}
}
