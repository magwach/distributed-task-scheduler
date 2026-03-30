package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type TaskHandlerImpl struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandlerImpl {
	return &TaskHandlerImpl{
		Service: service,
	}
}

func (h *TaskHandlerImpl) CreateTask(c *fiber.Ctx) error {

	task := dto.CreateTaskRequest{}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Invalid request body",
		})
	}

	if task.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Please provide a title",
		})
	}

	if task.Schedule == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Please provide a schedule",
		})
	}

	if task.Priority == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Please provide a priority",
		})
	}

	newTask, err := h.Service.CreateTask(task)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to create task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Task created successfully",
		"data":    newTask,
	})
}

func (h *TaskHandlerImpl) GetTasks(c *fiber.Ctx) error {

	tasks, err := h.Service.GetTasks()

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get tasks",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Tasks fetched successfully",
		"data":    tasks,
	})
}

func (h *TaskHandlerImpl) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse the id",
		})
	}

	task, err := h.Service.GetTask(parsedId)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Task fetched successfully",
		"data":    task,
	})
}

func (h *TaskHandlerImpl) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse the id",
		})
	}

	deletedTask, err := h.Service.DeleteTask(parsedId)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to delete task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Task deleted successfully",
		"data":    deletedTask,
	})
}

func (h *TaskHandlerImpl) GetLogs(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse the id",
		})
	}

	logs, err := h.Service.GetLogs(parsedId)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Logs fetched successfully",
		"data":    logs,
	})
}

func (h *TaskHandlerImpl) RetryTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse the id",
		})
	}
	err = h.Service.RetryTask(parsedId)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error":   err.Error(),
			"message": "Failed to retry task",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Task retried successfully",
	})
}
