package API

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		tasks []Task
		err   error
	)
	teamIDParam := c.Query("team_id")
	if teamIDParam == "" {
		tasks, err = DBAllTasks(ctx)
	} else {
		teamID, parseErr := strconv.Atoi(teamIDParam)
		if parseErr != nil || teamID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "team_id must be a positive integer"})
			return
		}
		tasks, err = DBTasksByTeam(ctx, teamID)
	}
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}

	fmt.Printf("📦 Получено задач: %d\n", len(tasks))
	fmt.Printf("📋 Данные: %+v\n", tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	ctx := c.Request.Context()

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task id must be a positive integer"})
		return
	}

	task, err := DBTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		log.Printf("failed to fetch task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task"})
		return
	}

	c.JSON(http.StatusOK, task)
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

func GetTeamByID(c *gin.Context) {
	ctx := c.Request.Context()

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil || teamID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team id must be a positive integer"})
		return
	}

	team, err := DBTeamByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		log.Printf("failed to fetch team: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch team"})
		return
	}

	c.JSON(http.StatusOK, team)
}
