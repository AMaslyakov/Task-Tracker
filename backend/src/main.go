package main

import (
	"backend/API"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var err error
	dsn := "postgres://shani:654321@localhost:5432/task-tracker?sslmode=disable"
	API.Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer API.Pool.Close()

	r := gin.Default()

	r.GET("/api/tasks", API.GetAllTasks)

	log.Println(" Сервер запущен на :8080")
	r.Run(":8080")
}
