package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/auth"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
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

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged in successfully",
	})
}

func (h *AuthHandlerImpl) Refresh(c *fiber.Ctx) error {
	return nil
}

func (h *AuthHandlerImpl) GoogleLogin(c *fiber.Ctx) error {
	state, err := auth.GenerateState()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to generate state",
		})
	}

	err = auth.StoreState(queue.RedisClient, state)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to store state",
		})
	}

	url := auth.GetGoogleAuthUrl(state)

	return c.Redirect(url)
}

func (h *AuthHandlerImpl) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing state or code"})
	}

	valid, err := auth.ValidateState(queue.RedisClient, state)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "state validation failed"})
	}

	if !valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid or expired state"})
	}

	user, err := auth.ExchangeGoogleCode(code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "oauth failed"})
	}

	newUser, err := h.Service.GetOrCreateOAuthUser(user.Email, user.Name, user.Picture, "google", user.ID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to login",
		})
	}

	token, err := auth.GenerateToken(newUser.ID, newUser.Email, newUser.Role)

	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to generate token",
		})
	}

	production := os.Getenv("ENVIROMENT") == "prod"

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   production,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Redirect("http://localhost:3000/dashboard", fiber.StatusTemporaryRedirect)
}

func (h *AuthHandlerImpl) GitHubLogin(c *fiber.Ctx) error {
	state, err := auth.GenerateState()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to generate state",
		})
	}

	err = auth.StoreState(queue.RedisClient, state)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to store state",
		})
	}

	url := auth.GetGithubAuthURL(state)

	return c.Redirect(url)
}

func (h *AuthHandlerImpl) GitHubCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing state or code"})
	}

	valid, err := auth.ValidateState(queue.RedisClient, state)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "state validation failed"})
	}

	if !valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid or expired state"})
	}

	user, err := auth.ExchangeGithubCode(code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "oauth failed"})
	}

	newUser, err := h.Service.GetOrCreateOAuthUser(user.Email, user.Name, user.AvatarURL, "github", strconv.Itoa(user.ID))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to login",
		})
	}

	token, err := auth.GenerateToken(newUser.ID, newUser.Email, newUser.Role)

	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"error": "failed to generate token",
		})
	}

	production := os.Getenv("ENVIROMENT") == "prod"

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   production,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Redirect("http://localhost:3000/dashboard", fiber.StatusTemporaryRedirect)
}
