package API

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

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

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task payload"})
		return
	}
	if err := validateUpdateTaskRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := DBUpdateTask(ctx, taskID, req)
	if err != nil {
		respondTaskDBError(c, err, "failed to update task")
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTaskStatus(c *gin.Context) {
	ctx := c.Request.Context()

	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task status payload"})
		return
	}
	req.StatusName = strings.TrimSpace(req.StatusName)
	if req.StatusName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status_name is required"})
		return
	}

	task, err := DBUpdateTaskStatus(ctx, taskID, req.StatusName)
	if err != nil {
		respondTaskDBError(c, err, "failed to update task status")
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := DBDeleteTask(ctx, taskID); err != nil {
		respondTaskDBError(c, err, "failed to delete task")
		return
	}

	c.Status(http.StatusNoContent)
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
