package main

import (
	"log"

	"github.com/eduardocaversan/lead-scraper-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/scrape", handler.ScrapeHandler)

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
