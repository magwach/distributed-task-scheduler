package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type AuthHandlerImpl struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		Service: service,
	}
}

func (h *AuthHandlerImpl) Register(c *fiber.Ctx) error {

	input := dto.UserRegister{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Email == "" || input.Name == "" || input.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Please provide all required fields",
		})
	}

	user, err := h.Service.Register(input)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "User registered successfully",
		"data":    user,
	})
}

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	input := dto.UserLogin{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Please provide all required fields",
		})
	}

	token, err := h.Service.Login(input)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to login",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"token": token,
	})
}

func (h *AuthHandlerImpl) Refresh(c *fiber.Ctx) error {
	return nil
}
