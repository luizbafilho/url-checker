package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Owner struct {
	Login string
}

type Repo struct {
	Name  string
	Owner Owner
}

// PullRequest represents a pull request struct
type PullRequest struct {
	Number int
	Body   string
	Repo   Repo
}

// URLStatus represents the status of a URL if it is reachable or not
type URLStatus struct {
	URL       string
	Reachable bool
}

type Checker interface {
	Check(pr *PullRequest) ([]URLStatus, error)
}

type UrlChecker struct {
	Timeout time.Duration
}

// very naive regex I got on google, it won't cover all cases, but the basic ones should
var reURL = regexp.MustCompile(`(?im)https?:\/\/[A-Z\d\.-]{2,}\.[A-Z]{2,}(:\d{2,4})?`)

// Check parses the pr's body looking for urls and checks their reachability
func (uc *UrlChecker) Check(pr *PullRequest) []URLStatus {
	// extracting all urls from pr's body
	matches := reURL.FindAllString(pr.Body, -1)
	matchesLen := len(matches)

	statusesCh := make(chan URLStatus, matchesLen)
	for _, m := range matches {
		go func(url string) {
			status, err := uc.urlReachable(url)
			if err != nil {
				log.Printf("url=%s is not reachable: %s", url, err)
			}
			statusesCh <- URLStatus{url, status}
		}(m)
	}

	statuses := []URLStatus{}
	for i := 0; i < matchesLen; i++ {
		statuses = append(statuses, <-statusesCh)
	}

	return statuses
}

func (uc *UrlChecker) urlReachable(u string) (bool, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false, err
	}

	address := parsedURL.Host
	if parsedURL.Port() == "" {
		switch parsedURL.Scheme {
		case "http":
			address += ":80"
		case "https":
			address += ":443"
		}
	}

	conn, err := net.DialTimeout("tcp", address, uc.Timeout)
	if err != nil {
		return false, err
	}
	conn.Close()

	return true, nil
}

type GHService struct {
	client *github.Client
}

func NewGHService(token string) *GHService {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GHService{
		client: client,
	}
}

// CreatePullRequestComment creates a new pull request comment with the reachability status of the urls
func (ghs *GHService) CreatePullRequestComment(pr *PullRequest, statuses []URLStatus) error {
	comment := "### Urls reachability report: \n"
	for _, s := range statuses {
		comment += fmt.Sprintf("- %s: %t\n", s.URL, s.Reachable)
	}
	input := &github.IssueComment{Body: &comment}

	_, _, err := ghs.client.Issues.CreateComment(context.Background(), pr.Repo.Owner.Login, pr.Repo.Name, pr.Number, input)
	if err != nil {
		return fmt.Errorf("Issues.CreateComment returned error: %v", err)
	}

	return nil
}

func main() {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	if githubAccessToken == "" {
		log.Fatal("Invalid GITHUB_ACCESS_TOKEN environment variable")
	}

	gsh := NewGHService(githubAccessToken)
	gsh.CreatePullRequestComment(&PullRequest{
		Number: 5,
		Repo: Repo{
			Name: "vimfiles",
			Owner: Owner{
				Login: "luizbafilho",
			},
		},
	}, []URLStatus{
		{"https://google.com", true},
		{"https://www.yahoo.com", false},
	})
}
