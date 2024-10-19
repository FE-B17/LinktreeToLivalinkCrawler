package main

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
)

// Result speichert die gecrawlten Informationen von der Linktree-Seite
type Result struct {
	Links       map[string]string `json:"links"`
	IconLinks   map[string]string `json:"icon_links"`
	Title       string            `json:"title"`
	ProfileName string            `json:"profile_name"`
	ProfileImg  string            `json:"profile_img"`
}

// crawlLinktreeProfile crawlt das angegebene Linktree-Profil und gibt die extrahierten Ergebnisse zurück
func crawlLinktreeProfile(link string) (Result, error) {
	collector := colly.NewCollector()

	var result Result
	result.Links = make(map[string]string)
	result.IconLinks = make(map[string]string)

	// Extrahiere den Seitentitel
	collector.OnHTML("head > title", func(e *colly.HTMLElement) {
		result.Title = e.Text
	})

	// Extrahiere reguläre Links
	collector.OnHTML("a[data-testid='LinkButton']", func(e *colly.HTMLElement) {
		linkText := e.ChildText("div > p")
		href := e.Attr("href")
		if linkText != "" && href != "" {
			result.Links[linkText] = href
		}
	})

	// Extrahiere soziale Icon-Links
	collector.OnHTML("a[data-testid='SocialIcon']", func(e *colly.HTMLElement) {
		iconName := e.ChildAttr("title", "title")
		href := e.Attr("href")
		if iconName != "" && href != "" {
			result.IconLinks[iconName] = href
		}
	})

	// Extrahiere den Profilnamen
	collector.OnHTML("div[id='profile-title']", func(e *colly.HTMLElement) {
		result.ProfileName = e.Text
	})

	// Extrahiere die Profilbild-URL
	collector.OnHTML("img[data-testid='ProfileImage']", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		if imgSrc != "" {
			result.ProfileImg = imgSrc
		}
	})

	// Logge die besuchte URL
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	// Besuche die Zielseite
	err := collector.Visit(link)
	if err != nil {
		return Result{}, fmt.Errorf("failed to visit %s: %w", link, err)
	}

	// Überprüfe, ob wichtige Felder gefunden wurden
	if result.Title == "" || result.ProfileName == "" {
		return Result{}, errors.New("failed to extract required profile information")
	}

	return result, nil
}
