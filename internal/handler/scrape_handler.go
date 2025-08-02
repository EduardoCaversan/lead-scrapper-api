package handler

import (
	"context"
	"log"
	"time"

	"github.com/eduardocaversan/lead-scraper-api/internal/model"
	"github.com/eduardocaversan/lead-scraper-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func ScrapeHandler(c *fiber.Ctx) error {
	var req model.ScrapeRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[Handler] Erro parsing JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao fazer parsing do corpo da requisição",
		})
	}

	if len(req.Keywords) == 0 {
		log.Printf("[Handler] Requisição sem keywords")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nenhuma palavra-chave fornecida",
		})
	}

	log.Printf("[Handler] Recebido scraping para keywords: %v", req.Keywords)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := service.ScrapeLeadsParallel(ctx, req.Keywords)
	if err != nil {
		log.Printf("[Handler] Erro durante scraping: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno ao realizar scraping",
		})
	}

	if results == nil {
		log.Printf("[Handler] Resultados vazios (nil)")
		return c.Status(fiber.StatusOK).JSON([]model.LeadResult{})
	}

	log.Printf("[Handler] Scraping concluído com %d resultados", len(results))
	return c.JSON(results)
}
