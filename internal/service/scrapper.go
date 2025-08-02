package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/eduardocaversan/lead-scraper-api/internal/model"
	"github.com/eduardocaversan/lead-scraper-api/internal/util"
)

const (
	duckduckgoHTML = "https://html.duckduckgo.com/html/?q=%s"
	maxConcurrent  = 2
)

func ScrapeLeadFromKeywordContext(ctx context.Context, keyword string) ([]model.LeadResult, error) {
	log.Printf("[Service] Iniciando scrape para keyword: %q", keyword)

	query := url.QueryEscape(keyword)
	searchURL := fmt.Sprintf(duckduckgoHTML, query)

	body, err := util.FetchHTML(ctx, searchURL)
	if err != nil {
		log.Printf("[Service] Erro FetchHTML para %s: %v", searchURL, err)
		return nil, err
	}

	log.Printf("[Service] HTML do DuckDuckGo obtido, tamanho: %d", len(body))

	results, err := parseDuckDuckGoResults(body, keyword)
	if err != nil {
		log.Printf("[Service] Erro parseDuckDuckGoResults: %v", err)
		return nil, err
	}

	log.Printf("[Service] Encontrados %d resultados para keyword %q", len(results), keyword)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var enrichedResults []model.LeadResult

	for _, res := range results {
		if ctx.Err() != nil {
			log.Printf("[Service] Contexto cancelado, abortando enrich de resultados")
			break
		}

		wg.Add(1)
		go func(r model.LeadResult) {
			defer wg.Done()

			log.Printf("[Service] Buscando contatos na página: %s", r.URL)
			contatos := scrapeContactsFromPage(ctx, r.URL)

			mu.Lock()
			if len(contatos.Emails) > 0 {
				r.Title += " | Emails: " + strings.Join(contatos.Emails, ", ")
			}
			if len(contatos.Phones) > 0 {
				r.Title += " | Telefones: " + strings.Join(contatos.Phones, ", ")
			}
			enrichedResults = append(enrichedResults, r)
			mu.Unlock()
		}(res)
	}

	wg.Wait()

	return removeDuplicateResults(enrichedResults), nil
}

type Contacts struct {
	Emails []string
	Phones []string
}

func scrapeContactsFromPage(ctx context.Context, pageURL string) Contacts {
	body, err := util.FetchHTML(ctx, pageURL)
	if err != nil {
		log.Printf("[Service] Erro ao buscar página %s: %v", pageURL, err)
		return Contacts{}
	}

	emails := util.ExtractEmails(body)
	phones := util.ExtractPhones(body)

	return Contacts{
		Emails: uniqueStrings(emails),
		Phones: uniqueStrings(phones),
	}
}

func uniqueStrings(input []string) []string {
	set := make(map[string]struct{}, len(input))
	var result []string
	for _, s := range input {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, exists := set[s]; !exists {
			set[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

// parseDuckDuckGoResults adapta a extração para links reais, incluindo decodificação de URLs relativas (uddg)
func parseDuckDuckGoResults(htmlBody, keyword string) ([]model.LeadResult, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear HTML: %w", err)
	}

	var results []model.LeadResult

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			class := getAttr(n, "class")
			if strings.Contains(class, "result__a") {
				rawLink := getAttr(n, "href")
				link := resolveDuckDuckGoLink(rawLink)
				title := extractText(n)

				if link != "" && !strings.HasPrefix(link, "#") && !strings.HasPrefix(link, "/") {
					results = append(results, model.LeadResult{
						Keyword: keyword,
						Title:   title,
						URL:     link,
					})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return results, nil
}

func resolveDuckDuckGoLink(raw string) string {
	if raw == "" {
		return ""
	}
	// Se o link começa com "/l/?kh=1&uddg=ENCODED_URL", decodifica para pegar o URL real
	if strings.HasPrefix(raw, "/l/?kh=") {
		u, err := url.Parse(raw)
		if err != nil {
			return raw
		}
		q := u.Query()
		uddg := q.Get("uddg")
		if uddg != "" {
			decoded, err := url.QueryUnescape(uddg)
			if err == nil {
				return decoded
			}
		}
		return raw
	}
	return raw
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func extractText(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractText(c)
	}
	return strings.TrimSpace(result)
}

func removeDuplicateResults(results []model.LeadResult) []model.LeadResult {
	seen := make(map[string]struct{})
	var unique []model.LeadResult
	for _, r := range results {
		if _, ok := seen[r.URL]; !ok {
			seen[r.URL] = struct{}{}
			unique = append(unique, r)
		}
	}
	return unique
}

func ScrapeLeadsParallel(ctx context.Context, keywords []string) ([]model.LeadResult, error) {
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allResults []model.LeadResult
	var firstErr error

	for _, keyword := range keywords {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sem <- struct{}{}:
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				defer func() { <-sem }()

				results, err := ScrapeLeadFromKeywordContext(ctx, k)
				if err != nil {
					log.Printf("[Service] Erro ao buscar keyword %q: %v", k, err)
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}

				mu.Lock()
				allResults = append(allResults, results...)
				mu.Unlock()
			}(keyword)
		}
	}

	wg.Wait()
	if allResults == nil {
		allResults = []model.LeadResult{}
	}

	return allResults, firstErr
}
