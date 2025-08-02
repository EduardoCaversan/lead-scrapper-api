package util

import (
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,}`)
	phoneRegex = regexp.MustCompile(`\(?\d{2}\)?[\s.-]?\d{4,5}[\s.-]?\d{4}`)
	urlRegex   = regexp.MustCompile(`https?://[^\s"]+`)
)

func ExtractEmails(html string) []string {
	return emailRegex.FindAllString(html, -1)
}

func ExtractPhones(html string) []string {
	return phoneRegex.FindAllString(html, -1)
}

func ExtractLinks(html string) []string {
	return urlRegex.FindAllString(html, -1)
}
