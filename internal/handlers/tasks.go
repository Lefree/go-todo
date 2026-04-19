package handlers

import (
	"lefree/go-todo/internal/models"
	"lefree/go-todo/internal/service"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type TaskHandler struct {
	logger  *slog.Logger
	service *service.TaskService
}

func NewTaskHandler(log *slog.Logger, s *service.TaskService) *TaskHandler {
	return &TaskHandler{
		logger:  log,
		service: s,
	}
}

func (h *TaskHandler) AddTask(c *echo.Context) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	id, err := h.service.AddTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]string{"result": id})
	return nil
}

func (h *TaskHandler) GetAllTasks(c *echo.Context) error {
	res, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
	return nil
}

func (h *TaskHandler) GetTaskById(c *echo.Context) error {
	taskId := c.Param("id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty param Id"})
	} else {
		response, err := h.service.GetById(taskId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if response.ID == "" {
			c.JSON(http.StatusNotFound, nil)
		}
		c.JSON(http.StatusOK, response)
	}
	return nil
}

func (h *TaskHandler) DeleteTaskById(c *echo.Context) error {
	taskId := c.Param("id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty param Id"})
	} else {
		err := h.service.Delete(taskId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		c.JSON(http.StatusOK, map[string]string{"result": "success"})
	}
	return nil
}
