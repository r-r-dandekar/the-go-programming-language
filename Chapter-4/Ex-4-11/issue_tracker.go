package main

import (
	"bufio"
	"fmt"
	"issues/github"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// TODO:
//  - Reading issues (in more detail?)
//  - The CLI

var client *github.Client
var appName = "GitHub Issue Tracker CLI"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to system environment variables")
	}
	token := os.Getenv("GH_ACCESS_TOKEN")

	client = github.NewClient(token)

	for true {
		reader := bufio.NewReader(os.Stdin)
		done := cliLoop(reader)
		if done {
			break
		}
	}
}

func cliLoop(reader *bufio.Reader) bool {
	fmt.Print("issue_tracker > ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v\n", err)
		return true
	}
	input = strings.TrimSpace(input)

	if str, match := TrimIfHasPrefix(input, "search "); match {
		cliSearchIssues(str)
	} else if str, match := TrimIfHasPrefix(input, "create "); match {
		cliCreateIssue(str)
	} else if _, match := TrimIfHasPrefix(input, "exit "); match {
		fmt.Printf("Exiting %v...\n", appName)
		return true
	}

	return false
}

func TrimIfHasPrefix(str, prefix string) (string, bool) {
	has := strings.HasPrefix(str, prefix)
	if has {
		return strings.TrimPrefix(str, prefix), true
	} else {
		return "", false
	}
}

func cliCreateIssue(str string) {
	args := strings.Fields(str)
	if len(args) < 2 {
		errorCreatingIssue()
		return
	}
	ownerRepo := strings.Split(args[0], "/")
	if len(ownerRepo) != 2 {
		errorCreatingIssue()
		return
	}
	owner, repo := ownerRepo[0], ownerRepo[1]
	title := args[1]

	var body string
	assignees := []string{"rrdan42"}
	labels := []string{"bug"}
	if len(args) > 2 {
		body = args[2]
	}

	issue := github.CreateIssueRequest{
		Title:     title,
		Body:      body,
		Assignees: assignees,
		Labels:    labels,
	}

	response, err := client.CreateIssue(owner, repo, issue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating issue: %v\n", err)
		return
	}
	fmt.Println("Response:\n", response)

}

func errorCreatingIssue() {
	fmt.Fprintf(os.Stderr, "Usage:\ncreate owner/repo whatever...")
}

func cliSearchIssues(str string) {
	keywords := strings.Fields(str)
	result, err := github.SearchIssues(keywords)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)
	oneYearAgo := now.AddDate(-1, 0, 0)

	fmt.Println("\n\tLess than a month ago:")
	flag := false
	for _, item := range result.Items {
		if item.CreatedAt.After(oneMonthAgo) {
			printIssue(item)
			flag = true
		}
	}
	if !flag {
		fmt.Println("\t\t< no issues >")
	}
	flag = false
	fmt.Println("\n\tLess than a year ago:")
	for _, item := range result.Items {
		if !item.CreatedAt.After(oneMonthAgo) && item.CreatedAt.After(oneYearAgo) {
			printIssue(item)
			flag = true
		}
	}
	if !flag {
		fmt.Println("\t\t< no issues >")
	}
	flag = false
	fmt.Println("\n\tMore than a year ago:")
	for _, item := range result.Items {
		if !item.CreatedAt.After(oneYearAgo) {
			printIssue(item)
			flag = true
		}
	}
	if !flag {
		fmt.Println("\t\t< no issues >")
	}
	flag = false
}
func printIssue(issue *github.Issue) {
	fmt.Printf("\t\t#%-5d %9.9s %.55s - %v\n",
		issue.Number, issue.User.Login, issue.Title, issue.CreatedAt)
}
