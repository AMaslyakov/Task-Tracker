package API

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := DBAllTasks(ctx)
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}

	fmt.Printf("📦 Получено задач: %d\n", len(tasks))
	fmt.Printf("📋 Данные: %+v\n", tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetAllTeams(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := DBAllTeams(ctx)
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}

	fmt.Printf("📦 Получено задач: %d\n", len(tasks))
	fmt.Printf("📋 Данные: %+v\n", tasks)
	c.JSON(http.StatusOK, tasks)
}
