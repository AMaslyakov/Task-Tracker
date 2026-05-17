package API

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteTask godoc
// @Summary Delete task
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/task/{id} [delete]
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

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = DBDeleteUser(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "пользователь удален"})
}
