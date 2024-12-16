package authhdl

import (
	"kn-assignment/internal/constant"
	errors "kn-assignment/internal/core/error"
	"kn-assignment/internal/handler/dto"
	"kn-assignment/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} errors.CustomError
// @Failure 500 {object} errors.CustomError
// @Router /auth/register [post]
func (h *handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var user dto.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewCustomError(constant.ErrCodeInvalidRequest))
		return
	}

	in := user.ToDomain()

	if err := h.svc.RegisterUser(ctx, in); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{Message: "User registered successfully"})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} errors.CustomError
// @Failure 401 {object} errors.CustomError
// @Failure 500 {object} errors.CustomError
// @Router /auth/login [post]
func (h *handler) Login(c *gin.Context) {
	var credentials dto.LoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewCustomError(constant.ErrCodeInvalidRequest))
		return
	}

	response, err := h.svc.AuthenticateUser(c.Request.Context(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	var loginResponse dto.LoginResponse

	c.JSON(http.StatusOK, loginResponse.FromDomain(response))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body dto.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} errors.CustomError
// @Failure 500 {object} errors.CustomError
// @Router /auth/refresh-token [post]
func (h *handler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewCustomError(constant.ErrCodeInvalidRequest))
		return
	}

	claims, err := util.ValidateToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewCustomError(constant.ErrCodeInvalidRequest))
		return
	}

	newAccessToken, err := util.GenerateAccessToken(claims.Id, claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewCustomError(constant.ErrCodeGenerateToken))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{AccessToken: newAccessToken})
}
