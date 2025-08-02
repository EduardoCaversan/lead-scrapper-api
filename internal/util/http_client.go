package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchHTML(ctx context.Context, url string) (string, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("falha ao criar request: %w", err)
	}

	// Defina User-Agent para evitar bloqueio por robô
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LeadScraperBot/1.0; +https://yourdomain.com/bot)")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("requisição retornou status %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler corpo da resposta: %w", err)
	}

	return string(body), nil
}
