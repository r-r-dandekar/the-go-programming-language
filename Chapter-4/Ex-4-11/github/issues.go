package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

const (
	LabelBug     string = "bug"
	LabelTask    string = "task"
	LabelFeature string = "feature"
)

type Client struct {
	token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{
		token:      token,
		HTTPClient: client,
	}
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := client.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	// We must close resp.Body on all execution paths.
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func (c *Client) GetIssue(owner, repo string, issueNumber int) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	tokenStr := fmt.Sprintf("Bearer %s", c.token)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api error (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	var issue GetIssueResponse
	if err := json.Unmarshal(respBodyBytes, &issue); err != nil {
		return "", fmt.Errorf("failed to decode response json: %w", err)
	}

	labels := []string{}
	for _, l := range issue.Labels {
		labels = append(labels, l.Name)
	}
	assignees := []string{}
	for _, u := range issue.Assignees {
		assignees = append(assignees, u.Login)
	}

	// pretty-printed result
	result := fmt.Sprintf("#%d - %s \t\t%v \n[%s]\n\n%v \t\tlast updated : %v\n"+
		"----------------------------------------------------------\n\n"+
		"%v\n\n"+
		"----------------------------------------------------------\n"+
		"Assigned to: %v", issue.Number, issue.Title, labels, issue.State, issue.User.Login, issue.UpdatedAt, issue.Body, assignees)

	return result, nil
}

func (c *Client) CreateIssue(owner, repo string, issue CreateIssueRequest) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

	jsonData, err := json.Marshal(issue)
	if err != nil {
		return "", fmt.Errorf("failed to marshal issue: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	tokenStr := fmt.Sprintf("Bearer %s", c.token)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("github api error (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	return string(respBodyBytes), nil
}

func (c *Client) UpdateIssue(owner, repo string, issueNumber int, issue map[string]any) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, issueNumber)

	jsonData, err := json.Marshal(issue)
	if err != nil {
		return "", fmt.Errorf("failed to marshal issue: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	tokenStr := fmt.Sprintf("Bearer %s", c.token)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api error (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	return string(respBodyBytes), nil
}
