// Issues prints a table of GitHub issues matching the search terms.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"issues/github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
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
