package main

import (
	"backend/API"
	"backend/events"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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

	brokerURL := os.Getenv("BROKER_URL")
	if brokerURL == "" {
		log.Fatal("BROKER_URL is required")
	}
	exchangeName := os.Getenv("BROKER_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "task_events"
	}

	ctx := context.Background()
	rmqClient, err := events.NewRabbitMQClient(ctx, brokerURL, exchangeName)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer rmqClient.Close()

	API.SetRabbitMQClient(rmqClient)
	log.Println("RabbitMQ connected")

	sseHub := events.NewSSEHub()
	go sseHub.Run()
	log.Println("SSE Hub started")
	go func() {
		if err := rmqClient.Consume(ctx, func(event events.Event) error {
			sseHub.Broadcast(event)
			return nil
		}); err != nil {
			log.Printf("RabbitMQ consumer stopped: %v", err)
		}
	}()
	log.Println("RabbitMQ consumer started")

	r := gin.Default()

	r.GET("/health", healthCheck)
	r.POST("/api/user", API.CreateUser)
	r.PATCH("/api/users/:id", API.UpdateUser)
	r.DELETE("/api/users/:id", API.DeleteUser)
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
	api.GET("/events", sseHub.SSEHandler)

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
