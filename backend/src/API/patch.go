package API

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UpdateTask godoc
// @Summary Update task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param payload body UpdateTaskRequest true "Task update payload"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/task/{id} [patch]
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

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Changes task status, intended for drag-and-drop between dashboard columns.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param payload body UpdateTaskStatusRequest true "Task status payload"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/task/{id}/status [patch]
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

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json payload"})
		return
	}

	err = DBUpdateUser(c.Request.Context(), id, req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "payload must contain at least one field":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}
