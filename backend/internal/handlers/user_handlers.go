package handlers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type UserHandlerImpl struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		Service: service,
	}
}

func (h *UserHandlerImpl) Me(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	user, err := h.Service.GetUser(email)

	if err != nil {
		log.Println("Failed to get user: ", err)
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":    userID,
		"email":      email,
		"role":       role,
		"avatar_url": user.AvatarUrl,
	})
}
