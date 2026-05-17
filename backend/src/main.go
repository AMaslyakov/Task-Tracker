package main

import (
	"backend/API"
	_ "backend/docs"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Tracker API
// @version 1.0
// @description Backend API for Victory Group Task Tracker.
// @host localhost:8081
// @BasePath /
func main() {
	var err error
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	API.Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer API.Pool.Close()

	r := gin.Default()

	r.GET("/health", healthCheck)
	r.GET("/api/tasks", API.GetAllTasks)
	r.GET("/api/task/:id", API.GetTaskByID)
	r.POST("/api/task", API.InsertTask)
	r.POST("/api/user", API.CreateUser)
	r.PATCH("/api/task/:id", API.UpdateTask)
	r.PATCH("/api/task/:id/status", API.UpdateTaskStatus)
	r.PATCH("/api/users/:id", API.UpdateUser)
	r.DELETE("/api/task/:id", API.DeleteTask)
	r.DELETE("/api/users/:id", API.DeleteUser)
	r.GET("/api/teams", API.GetAllTeams)
	r.GET("/api/team/:id", API.GetTeamByID)
	r.POST("/api/login", API.Login)
	r.GET("/api/me", API.Me)
	r.POST("/api/logout", API.Logout)

	api := r.Group("/api")
	api.Use(API.RequireAuth())
	api.GET("/tasks", API.GetAllTasks)
	api.GET("/task/:id", API.GetTaskByID)
	api.POST("/task", API.InsertTask)
	api.PATCH("/task/:id", API.UpdateTask)
	api.PATCH("/task/:id/status", API.UpdateTaskStatus)
	api.DELETE("/task/:id", API.DeleteTask)
	api.GET("/teams", API.GetAllTeams)
	api.GET("/team/:id", API.GetTeamByID)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println(" Сервер запущен на :8080")
	r.Run(":8080")
}

// healthCheck godoc
// @Summary Check backend health
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
