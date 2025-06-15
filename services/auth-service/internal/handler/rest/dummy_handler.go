package rest

import (
	"github.com/content-management-system/auth-service/internal/model/types"
	"github.com/content-management-system/auth-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	userService *service.UserService
	logger      *logrus.Logger
}

func NewHandler(userService *service.UserService, logger *logrus.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  types.User `json:"user"`
}

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.ValidatePassword(req.Email, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("Login failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token here (implement your JWT logic)
	token := "your-jwt-token-here" // Replace with actual JWT generation

	// Don't send password in response
	user.Password = ""

	return c.JSON(AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("Registration failed")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	token := "rhIrItXMIMoWbPYAywEJiwrSHlx/hBad9LeErTyfFBmmAus+yAp6OHu2NhojFxp06oWGEoJrLWomvir+Boc=\n"
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) GetUsersHandler(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get users")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

func (h *Handler) GetUserHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	user.Password = ""

	return c.JSON(fiber.Map{
		"data": user,
	})
}
