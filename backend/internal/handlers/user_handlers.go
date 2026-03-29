package handlers

import "github.com/gofiber/fiber/v2"

type UserHandlerImpl struct {
}

func NewUserHandler() *UserHandlerImpl {
	return &UserHandlerImpl{}
}

func (h *UserHandlerImpl) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	email := c.Locals("email").(string)
	role := c.Locals("role").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}
