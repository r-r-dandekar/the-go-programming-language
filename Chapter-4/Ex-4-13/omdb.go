package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var key string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to system environment variables")
	}
	key = os.Getenv("OMDB_API_KEY")
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("Enter movie title: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		title := strings.TrimSpace(line)

		err = searchByTitle(title)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error downloading poster:", err)
		} else {
			fmt.Println("Saved poster")
		}
	}
}

func searchByIMDb(id string) error {
	reqURL := fmt.Sprintf("https://img.omdbapi.com/?i=%s&apikey=%s", url.QueryEscape(id), url.QueryEscape(key))
	filename := id + ".jpeg"
	err := downloadFile(reqURL, filename)
	return err
}

func searchByTitle(title string) error {
	reqURL := fmt.Sprintf("https://www.omdbapi.com/?t=%s&apikey=%s", url.QueryEscape(title), url.QueryEscape(key))
	resp, err := http.Get(reqURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var res struct {
		ImdbID string `json:"imdbID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return searchByIMDb(res.ImdbID)
}

func downloadFile(rawURL, destPath string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return fmt.Errorf("request failed for %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found (404): %s", rawURL)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected HTTP status %d for %s", resp.StatusCode, rawURL)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write response to %s: %w", destPath, err)
	}

	return nil
}
