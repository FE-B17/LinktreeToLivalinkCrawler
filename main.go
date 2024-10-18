package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Result stores the crawled information from the Linktree page
type Result struct {
	Links       map[string]string `json:"links"`
	IconLinks   map[string]string `json:"icon_links"`
	Title       string            `json:"title"`
	ProfileName string            `json:"profile_name"`
	ProfileImg  string            `json:"profile_img"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [profile_name]")
		return
	}

	profileName := os.Args[1]
	link := fmt.Sprintf("https://linktr.ee/%s", profileName)

	result, err := crawlLinktreeProfile(link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during crawling: %v\n", err)
		return
	}

	// Save the crawled data to [ProfileName].json
	if err := saveResultToFile(result, result.ProfileName+".json"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save result: %v\n", err)
		return
	}

	// Download and save the profile image
	if result.ProfileImg != "" {
		if err := downloadProfileImage(result.ProfileImg, profileName+".jpg"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to download profile image: %v\n", err)
		}
	}
}

// crawlLinktreeProfile crawls the specified Linktree profile and returns the extracted result.
func crawlLinktreeProfile(link string) (Result, error) {
	collector := colly.NewCollector()

	var result Result
	result.Links = make(map[string]string)
	result.IconLinks = make(map[string]string)

	// Extract profile title
	collector.OnHTML("head > title", func(e *colly.HTMLElement) {
		result.Title = e.Text
	})

	// Extract regular links
	collector.OnHTML("a[data-testid='LinkButton']", func(e *colly.HTMLElement) {
		linkText := e.ChildText("div > p")
		href := e.Attr("href")
		if linkText != "" && href != "" {
			result.Links[linkText] = href
		}
	})

	// Extract social icon links
	collector.OnHTML("a[data-testid='SocialIcon']", func(e *colly.HTMLElement) {
		iconName := e.ChildAttr("title", "title")
		href := e.Attr("href")
		if iconName != "" && href != "" {
			result.IconLinks[iconName] = href
		}
	})

	// Extract profile name
	collector.OnHTML("div[id='profile-title']", func(e *colly.HTMLElement) {
		result.ProfileName = e.Text
	})

	// Extract profile image URL
	collector.OnHTML("img[data-testid=\"ProfileImage\"]", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		if imgSrc != "" {
			result.ProfileImg = imgSrc
		}
	})

	// Error logging during requests
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	// Visit the target profile URL
	err := collector.Visit(link)
	if err != nil {
		return Result{}, fmt.Errorf("failed to visit %s: %w", link, err)
	}

	collector.Wait()

	// Ensure that mandatory fields are captured
	if result.Title == "" || result.ProfileName == "" {
		return Result{}, errors.New("failed to extract required profile information")
	}

	return result, nil
}

// saveResultToFile saves the crawled result to a JSON file.
func saveResultToFile(result Result, fileName string) error {
	// Create output path
	outputPath := filepath.Join(".", fileName)

	// Convert the result to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling result to JSON: %w", err)
	}

	// Write the JSON data to a file
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing result to file: %w", err)
	}

	fmt.Println("Data successfully saved to", outputPath)
	return nil
}

// downloadProfileImage downloads the profile image from the provided URL and saves it locally.
func downloadProfileImage(imgURL, fileName string) error {
	// Get the data from the URL
	response, err := http.Get(imgURL)
	if err != nil {
		return fmt.Errorf("failed to download image from %s: %w", imgURL, err)
	}
	defer response.Body.Close()

	// Check if the request was successful
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, status code: %d", response.StatusCode)
	}

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to save image to file: %w", err)
	}

	fmt.Println("Profile image successfully saved as", fileName)
	return nil
}
