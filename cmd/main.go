package main

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"log"
	"trace-example/server"
	"trace-example/storage"
)

func main() {
	app := fiber.New()

	// Подключаемся к Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.Fatal("create redis client", err)
	}

	// Настраиваем роутер
	handler := server.NewFiberHandler(storage.NewNotesStorage(client))
	app.Post("/create", handler.CreateNote)
	app.Get("/get", handler.GetNote)

	log.Fatal(app.Listen(":8080"))
}
