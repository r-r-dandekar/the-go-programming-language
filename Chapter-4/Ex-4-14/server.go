// A web server to search for GitHub issues and navigate through them
package main

import (
	"fmt"
	"html/template"
	"issues/github"
	"log"
	"net/http"
	"strconv"
)

var issueList *github.IssuesSearchResult

func main() {
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/results", resultsPage)
	http.HandleFunc("/issue", issuePage)
	address := "localhost:8000"
	fmt.Printf("Starting server at: http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	if err := landingTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}
func resultsPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	if len(q) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var err error
	issueList, err = github.SearchIssues(q)
	if err != nil {
		log.Fatal(err)
	}
	if err := issueListTemplate.Execute(w, issueList); err != nil {
		log.Fatal(err)
	}
}
func issuePage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")
	id, err := strconv.Atoi(q)
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't convert to id integer: %w", err))
	}
	if issueList == nil {
		log.Fatal(fmt.Errorf("no issues found!"))
	}
	if id >= len(issueList.Items) {
		log.Fatal(fmt.Errorf("id %d refers to an invalid issue!", id))
	}
	if err := issueTemplate.Execute(w, issueList.Items[id]); err != nil {
		log.Fatal(err)
	}
}

var landingTemplate = template.Must(template.ParseFiles("templates/index.html"))
var issueTemplate = template.Must(template.ParseFiles("templates/issue.html"))
var issueListTemplate = template.Must(template.ParseFiles("templates/issue_list.html"))
