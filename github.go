package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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
func (ghs *GHService) CreatePullRequestComment(pr *PullRequest, repo *Repo, statuses []URLStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	comment := "### Urls reachability report: \n"
	for _, s := range statuses {
		comment += fmt.Sprintf("- %s: %t\n", s.URL, s.Reachable)
	}
	input := &github.IssueComment{Body: &comment}

	_, _, err := ghs.client.Issues.CreateComment(context.Background(), repo.Owner.Login, repo.Name, pr.Number, input)
	if err != nil {
		return fmt.Errorf("Issues.CreateComment returned error: %v", err)
	}

	log.Printf("URL reachability report succesfully created for Pull Request: %d in %s repository from %s", pr.Number, repo.Name, repo.Owner.Login)

	return nil
}
