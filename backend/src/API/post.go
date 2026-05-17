package API

import (
	"backend/events"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// InsertTask godoc
// @Summary Create task
// @Tags tasks
// @Accept json
// @Produce json
// @Param payload body CreateTaskRequest true "Task payload"
// @Success 201 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/task [post]
func InsertTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task payload"})
		return
	}
	if err := validateCreateTaskRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := DBInsertTask(ctx, req)
	if err != nil {
		respondTaskDBError(c, err, "failed to create task")
		return
	}

	payload := events.TaskPayload{
		TaskID: task.ID,
		TeamID: task.TeamID,
		Title:  task.Title,
		Status: task.StatusName,
	}
	go publishEvent(events.NewEvent(events.EventTaskCreated, payload))

	c.JSON(http.StatusCreated, task)
}

func parseTaskID(c *gin.Context) (int, bool) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task id must be a positive integer"})
		return 0, false
	}
	return taskID, true
}

func validateCreateTaskRequest(req *CreateTaskRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	req.StatusName = strings.TrimSpace(req.StatusName)

	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.StatusName == "" {
		return errors.New("status_name is required")
	}
	if req.PriorityID <= 0 {
		return errors.New("priority_id must be a positive integer")
	}
	if req.TeamID <= 0 {
		return errors.New("team_id must be a positive integer")
	}
	if req.CreatedBy <= 0 {
		return errors.New("created_by must be a positive integer")
	}
	if req.AssignedTo != nil && *req.AssignedTo <= 0 {
		return errors.New("assigned_to must be a positive integer")
	}
	return nil
}

func validateUpdateTaskRequest(req *UpdateTaskRequest) error {
	hasField := false

	if req.Title != nil {
		hasField = true
		trimmed := strings.TrimSpace(*req.Title)
		req.Title = &trimmed
		if trimmed == "" {
			return errors.New("title must not be empty")
		}
	}
	if req.Description != nil {
		hasField = true
	}
	if req.StatusName != nil {
		hasField = true
		trimmed := strings.TrimSpace(*req.StatusName)
		req.StatusName = &trimmed
		if trimmed == "" {
			return errors.New("status_name must not be empty")
		}
	}
	if req.PriorityID != nil {
		hasField = true
		if *req.PriorityID <= 0 {
			return errors.New("priority_id must be a positive integer")
		}
	}
	if req.Deadline != nil {
		hasField = true
	}
	if req.ClearDeadline {
		hasField = true
	}
	if req.TeamID != nil {
		hasField = true
		if *req.TeamID <= 0 {
			return errors.New("team_id must be a positive integer")
		}
	}
	if req.AssignedTo != nil {
		hasField = true
		if *req.AssignedTo <= 0 {
			return errors.New("assigned_to must be a positive integer")
		}
	}
	if req.ClearAssignedTo {
		hasField = true
	}
	if !hasField {
		return errors.New("task update payload must contain at least one field")
	}
	return nil
}

func respondTaskDBError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if errors.Is(err, ErrRelatedEntityNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "related entity not found"})
		return
	}

	log.Printf("%s: %v", fallbackMessage, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMessage})
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	newID, err := DBCreateUser(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "пользователь с таким именем или email уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"id":      newID,
	})
}
