package main

import (
	"database/sql"
	"fmt"
	"lefree/go-todo/internal/handlers"
	"lefree/go-todo/internal/service"
	"log/slog"

	"github.com/labstack/echo/v5"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "todo"
)

func main() {
	e := echo.New()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})
	logger := slog.Default()
	taskService := service.NewTaskService(logger, db)
	handler := handlers.NewTaskHandler(logger, taskService)
	e.GET("/", handler.GetAllTasks)
	e.GET("/:id", handler.GetTaskById)
	e.POST("/", handler.AddTask)
	e.DELETE("/:id", handler.DeleteTaskById)

	e.Start(":9091")
}
