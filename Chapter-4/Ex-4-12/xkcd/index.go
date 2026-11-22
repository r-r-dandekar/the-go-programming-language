package xkcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: review error messages

type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Dau        string `json:"day"`
}

type ComicsFile struct {
	Comics []Comic `json:"comics"`
}

// store: appends the entry to the relevant file without checking if it already exists
func store(c *Comic) error {
	var data ComicsFile

	filename := indexFileName(c.Num)

	// 1. Read the existing file if it exists
	raw, err := os.ReadFile(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}
		// File doesn't exist yet; start with an empty slice so it serializes as [] instead of null
		data.Comics = []Comic{}
	} else if len(raw) > 0 {
		// 2. Unmarshal existing JSON into the wrapper struct
		if err := json.Unmarshal(raw, &data); err != nil {
			return fmt.Errorf("failed to parse JSON in %s: %w", filename, err)
		}
		if data.Comics == nil {
			data.Comics = []Comic{}
		}
	} else {
		data.Comics = []Comic{}
	}

	// 3. Append the new struct to the "comics" slice
	data.Comics = append(data.Comics, *c)

	// 4. Convert back to formatted JSON
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	updatedJSON = append(updatedJSON, '\n')

	// 5. Write back to the file
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filename, err)
	}
	if err := os.WriteFile(filename, updatedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %w", filename, err)
	}

	return nil
}

func Exists(n int) (bool, error) {
	filename := indexFileName(n)

	raw, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	if len(raw) == 0 {
		return false, nil
	}

	var data ComicsFile
	if err := json.Unmarshal(raw, &data); err != nil {
		return false, fmt.Errorf("failed to parse JSON in %s: %w", filename, err)
	}

	for _, comic := range data.Comics {
		if comic.Num == n {
			return true, nil
		}
	}

	return false, nil
}

// Inserts a Comic into the index if it's not already there
// Returns (true, nil) if inserted, and (false, nil) if the comic was already present
func Put(c *Comic) (bool, error) {
	present, err := Exists(c.Num)
	if err != nil {
		return false, err
	}
	if present {
		return false, nil
	} else {
		err := store(c)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

// Reads a comic from the index
func Get(n int) (*Comic, error) {
	present, err := Exists(n)
	if err != nil {
		return nil, err
	}
	if !present {
		return FetchAndPut(n)
	}

	filename := indexFileName(n)

	raw, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("index doesn't exist %s: %w", filename, err)
		}
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("comic not found in %s", filename)
	}

	var data ComicsFile
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s: %w", filename, err)
	}

	for _, comic := range data.Comics {
		if comic.Num == n {
			return &comic, nil
		}
	}
	return nil, nil
}

// Fetches a comic from xkcd.com
func fetch(n int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get comic %d: %s", n, resp.Status)
	}
	var c Comic
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func indexFileName(n int) string {
	return fmt.Sprintf("data/%dxx.json", (n / 100))
}

func FetchAndPut(n int) (*Comic, error) {
	comic, err := fetch(n)
	if err != nil {
		return nil, err
	}

	result, err := Put(comic)
	if err != nil {
		return comic, err
	}
	if result == false {
		return comic, nil
	}
	return comic, nil
}
