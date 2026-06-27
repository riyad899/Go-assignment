package handler

import (
	"net/http"

	"gotickets/internal/dto"
	"gotickets/internal/service"
	"gotickets/internal/utils"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Registration failed", err.Error())
	}

	res := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return utils.SendSuccess(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		return utils.SendError(c, http.StatusUnauthorized, "Login failed", err.Error())
	}

	res := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return utils.SendSuccess(c, http.StatusOK, "Login successful", res)
}
