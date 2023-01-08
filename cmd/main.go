package main

import (
	"context"
	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"log"
	"trace-example/server"
	"trace-example/storage"
	"trace-example/trace"
)

func main() {
	app := fiber.New()
	app.Use(otelfiber.Middleware("my-server"))

	tracer, err := trace.InitTracer("http://localhost:14268/api/traces", "Note Service")
	if err != nil {
		log.Fatal("init tracer", err)
	}

	// Подключаемся к Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.Fatal("create redis client", err)
	}
	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Fatal("enable instrument tracing", err)
	}

	// Настраиваем роутер
	handler := server.NewFiberHandler(storage.NewNotesStorage(client), tracer)
	app.Post("/create", handler.CreateNote)
	app.Get("/get", handler.GetNote)

	log.Fatal(app.Listen(":8080"))
}
