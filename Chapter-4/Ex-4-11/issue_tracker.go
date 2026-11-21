package main

import (
	"bufio"
	"flag"
	"fmt"
	"issues/github"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buildkite/shellwords"
	"github.com/joho/godotenv"
)

// TODO:
//  - Polishing
//		- Closing with variations
//		- More update options
//		- Support for multiple labels and assignees

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
	} else if str, match := TrimIfHasPrefix(input, "get "); match {
		cliGetIssue(str)
	} else if str, match := TrimIfHasPrefix(input, "update "); match {
		cliUpdateIssue(str)
	} else if str, match := TrimIfHasPrefix(input, "close "); match {
		cliCloseIssue(str)
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
func cliGetIssue(str string) {
	args, err := shellwords.Split(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading args while getting issue: %v\n\n", err)
		printGetUsage()
		return
	}

	if len(args) < 1 {
		printGetUsage()
		return
	}
	ownerRepoNumber := strings.Split(args[0], "#")
	if len(ownerRepoNumber) != 2 {
		printGetUsage()
		return
	}
	ownerRepo := strings.Split(ownerRepoNumber[0], "/")
	if len(ownerRepo) != 2 {
		printGetUsage()
		return
	}
	owner, repo, issueNumberStr := ownerRepo[0], ownerRepo[1], ownerRepoNumber[1]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse issue number as int while getting issue: %v\n\n", err)
		printGetUsage()
		return
	}

	// Pretty-printing of the result is handled in client.GetIssue
	result, err := client.GetIssue(owner, repo, issueNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting issue: %v\n", err)
		return
	}
	fmt.Printf("\n%s\n\n", result)
}

func printGetUsage() {
	fmt.Fprintln(os.Stderr, "Usage:\nget owner/repo#<issue-number>")
}

func cliCloseIssue(str string) {
	args, err := shellwords.Split(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading args while closing issue: %v\n\n", err)
		printCloseUsage()
		return
	}

	if len(args) < 1 {
		printCloseUsage()
		return
	}
	ownerRepoNumber := strings.Split(args[0], "#")
	if len(ownerRepoNumber) != 2 {
		printCloseUsage()
		return
	}
	ownerRepo := strings.Split(ownerRepoNumber[0], "/")
	if len(ownerRepo) != 2 {
		printCloseUsage()
		return
	}
	owner, repo, issueNumberStr := ownerRepo[0], ownerRepo[1], ownerRepoNumber[1]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse issue number as int while closing issue: %v\n\n", err)
		printCloseUsage()
		return
	}

	issue := make(map[string]any)
	issue["state"] = "closed"
	_, err = client.UpdateIssue(owner, repo, issueNumber, issue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting issue: %v\n", err)
		return
	} else {
		fmt.Printf("Issue %s/%s#%d closed\n", owner, repo, issueNumber)
	}
}

func printCloseUsage() {
	fmt.Fprintln(os.Stderr, "Usage:\nclose owner/repo#<issue-number>")
}

func cliCreateIssue(str string) {
	args, err := shellwords.Split(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading args while creating issue: %v\n\n", err)
		printCreateUsage()
		return
	}

	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	assignee := fs.String("a", "rrdan42", "issue assignees")
	label := fs.String("l", "bug", "issue labels")
	editor := fs.Bool("e", false, "use an external text editor for the body")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing args while creating issue: %v\n\n", err)
		printCreateUsage()
		return
	}

	var body string
	assignees := []string{*assignee}
	labels := []string{*label}
	positionalArgs := fs.Args()
	if len(positionalArgs) < 2 {
		printCreateUsage()
		return
	}
	ownerRepo := strings.Split(positionalArgs[0], "/")
	if len(ownerRepo) != 2 {
		printCreateUsage()
		return
	}
	owner, repo := ownerRepo[0], ownerRepo[1]
	title := positionalArgs[1]

	if *editor {
		var err error
		body, err = GetInputFromEditor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting input from editor while creating issue: %v\n\n", err)
			printCreateUsage()
			return
		}
	} else if len(positionalArgs) > 2 {
		body = positionalArgs[2]
	}

	issue := github.CreateIssueRequest{
		Title:     title,
		Body:      body,
		Assignees: assignees,
		Labels:    labels,
	}

	_, err = client.CreateIssue(owner, repo, issue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating issue: %v\n", err)
		return
	} else {
		fmt.Printf("Issue created\n")
	}
}

func printCreateUsage() {
	fmt.Fprintln(os.Stderr, "Usage:\ncreate [-l=label] [-a=assignee] [-e] owner/repo title [body]")
}
func cliUpdateIssue(str string) {
	args, err := shellwords.Split(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading args while creating issue: %v\n\n", err)
		printUpdateUsage()
		return
	}

	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	assignee := fs.String("a", "", "issue assignees")
	label := fs.String("l", "", "issue labels")
	title := fs.String("t", "", "issue title")
	body := fs.String("b", "", "issue body")
	editor := fs.Bool("e", false, "use an external text editor for the body")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing args while creating issue: %v\n\n", err)
		printUpdateUsage()
		return
	}

	setFlags := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	issue := make(map[string]any)
	if setFlags["a"] {
		issue["assignees"] = []string{*assignee}
	}
	if setFlags["l"] {
		issue["labels"] = []string{*label}
	}
	if setFlags["t"] {
		issue["title"] = title
	}
	if setFlags["b"] {
		issue["body"] = body
	}

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		printUpdateUsage()
		return
	}
	ownerRepoNumber := strings.Split(positionalArgs[0], "#")
	if len(ownerRepoNumber) != 2 {
		printUpdateUsage()
		return
	}
	ownerRepo := strings.Split(ownerRepoNumber[0], "/")
	if len(ownerRepo) != 2 {
		printUpdateUsage()
		return
	}
	owner, repo, issueNumberStr := ownerRepo[0], ownerRepo[1], ownerRepoNumber[1]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse issue number as int while getting issue: %v\n\n", err)
		printGetUsage()
		return
	}

	if *editor {
		var err error
		issue["body"], err = GetInputFromEditor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting input from editor while creating issue: %v\n\n", err)
			printUpdateUsage()
			return
		}
	}

	_, err = client.UpdateIssue(owner, repo, issueNumber, issue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating issue: %v\n", err)
		return
	} else {
		fmt.Printf("Issue %s/%s#%d updated\n", owner, repo, issueNumber)
	}
}

func printUpdateUsage() {
	fmt.Fprintln(os.Stderr, "Usage:\nupdate [-l=label] [-a=assignee] [-t=title] [-b=body] [-e] owner/repo#<issue-number>")
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
