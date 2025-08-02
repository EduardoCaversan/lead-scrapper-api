package model

type ScrapeRequest struct {
	Keywords []string `json:"keywords"`
}

type LeadResult struct {
	Keyword string `json:"keyword"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}
