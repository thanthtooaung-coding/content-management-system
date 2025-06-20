package auth_service

import (
	"github.com/content-management-system/auth-service/internal/service"
	"github.com/content-management-system/auth-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(us *service.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleID   uint64 `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password, req.RoleID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not generate refresh token")
	}

	return c.JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid refresh token request")
	}

	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	newToken, err := utils.GenerateToken(claims.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate new token")
	}

	return c.JSON(fiber.Map{"access_token": newToken})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
